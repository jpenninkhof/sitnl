'use strict';

var header = new Vue({
  el: '#header',
  data() {
    return {
      remainingDays: null,
      remainingHours: null,
    };
  },
  mounted() {
    this.getRemainingDates();
  },
  methods: {
    getRemainingDates() {
      const targetDate = new Date("2025-06-05T09:00:00.000+02:00").getTime();
      const currentDate = new Date().getTime();
      const timeRemaining = targetDate - currentDate;
      this.remainingDays = Math.floor(timeRemaining / (1000 * 60 * 60 * 24));
      this.remainingHours = Math.floor((timeRemaining % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    }
  },
});

var main = new Vue({
  el: '#main',
  data() {
    return {
      isCalendarsVisible: false,
      isExpanded: false,
      isModalVisible: false,
      showCalBlockers: false,
      eventTime: "event",
      filter: "all",
      activeSpeakers: null,
      lineup: [],
      proposalLineup: [],
      formattedLineup: [],
      sponsorLineup: [],
      expertCornerLineupUnsorted: [],
      expertCornerLineup: {},
      lastFocussedElementID: '',
      proposalLineupJson: proposalLineupJson,
      sponsorLineupJson: sponsorLineupJson
    };
  },
  mounted() {
    this.lineup = proposalLineupJson;
    this.formattedLineup = this.formatLineup();
    this.sponsorLineup = this.formatSponsorLineup();
    this.groupExpertCornerTopics();

    this.updateLiveSession();

    let interval;

    let timeNow = new Date().toISOString();

    const startCounterTime = new Date(
      "2025-06-05T00:50:00.000+02:00"
    ).toISOString();

    const endCounterTime = new Date(
      "2025-06-05T18:10:00.000+02:00"
    ).toISOString();

    if (timeNow > startCounterTime && timeNow <= endCounterTime) {
      interval = setInterval(() => {
        timeNow = new Date().toISOString();
        if (timeNow > endCounterTime) {
          clearInterval(interval);
          return;
        }
        this.updateLiveSession();
      }, 30000);
    }
  },
  methods: {
    toggleCalendars() {
      this.isCalendarsVisible = !this.isCalendarsVisible;
    },
    showSubscribeModal() {
      this.isModalVisible = !this.isModalVisible;
    },
    onFilterChange($event) {
      this.filter = $event.target.value;
      this.formattedLineup = this.formatLineup();
    },
    onTimeChange($event) {
      this.eventTime = $event.target.value;
    },
    getLocalTimeZone() {
      return luxon.DateTime.now().toFormat("Z");
    },
    formatLineup() {
      const tempLineUp = this.lineup.map((session) => {
        session.speakers.map((speaker) => {
          if(speaker.blueskyHandle) {
            speaker.blueskyHandle = this.formatBlueskyLink(speaker.blueskyHandle);
          }

          if(speaker.linkedInUrl) {
            speaker.linkedInUrl = this.formatLinkedInLink(speaker.linkedInUrl);
          }

          if(speaker.mastodonHandle) {
            speaker.mastodonHandle = this.formatMastodonLink(speaker.mastodonHandle);
          }

          if(speaker.blueskyHandle) {
            speaker.blueskyHandle = this.formatBlueskyLink(speaker.blueskyHandle);
          }
        });

        let start = session.startTime;
        let end = session.endTime;

        if (start === undefined) start = "00:00";
        if (end === undefined) end = "00:00";

        let tempStart = start.substring(0, start.indexOf(":"));
        let tempEnd = end.substring(0, end.indexOf(":"));

        if (tempStart.length == 1 && !tempStart.startsWith("0")) {
          start = "0" + start;
        }

        if (tempEnd.length == 1 && !tempEnd.startsWith("0")) {
          end = "0" + end;
        }

        let newStartTime = "2025-11-29T" + start + ":00.000+02:00";
        let newEndTime = "2025-11-29T" + end + ":00.000+02:00";

        if (session.title === "Lunch") {
          newStartTime = "2025-11-29T" + start + ":30.000+02:00";
        }

        let timeNow = new Date().toISOString();
        let sessionTimeStart = new Date(newStartTime).toISOString();
        let sessionTimeEnd = new Date(newEndTime).toISOString();
        let sessionLiveStatus = false;

        if (timeNow > sessionTimeStart && timeNow < sessionTimeEnd) {
          sessionLiveStatus = true;
        }

        return {
          ...session,
          startTime: newStartTime,
          endTime: newEndTime,
          isLive: sessionLiveStatus
        };
      });

      const sortedScheduleTemp = tempLineUp.sort((a, b) => {
        if (a.location === undefined) a.location = "";
        return luxon.DateTime.fromISO(a.startTime) -
        luxon.DateTime.fromISO(b.startTime) || a.location.localeCompare(b.location)
      });

      this.expertCornerLineupUnsorted = sortedScheduleTemp.filter(
        (schedule) => schedule.type.includes("expert")
      );

      const sortedSchedule = sortedScheduleTemp.filter(
        (schedule) => !schedule.type.includes("expert")
      );

      if (this.filter === "all") {
        return sortedSchedule;
      } else if (this.filter === "presentation_tech") {
        return sortedSchedule.filter((schedule) =>
          schedule.type.includes("presentation_tech")
        );
      } else if (this.filter === "presentation_bi") {
        return sortedSchedule.filter((schedule) =>
          schedule.type.includes("presentation_bi")
        );
      } else if (this.filter === "networking") {
        return sortedSchedule.filter((schedule) =>
          schedule.location.includes("networking")
        );
      } else if (this.filter === "creative") {
        return sortedSchedule.filter((schedule) =>
          schedule.location.includes("creative")
        );
      } else if (this.filter === "universe") {
        return sortedSchedule.filter((schedule) =>
          schedule.location.includes("universe")
        );
      } else {
        return sortedSchedule;
      }
    },
    formatSponsorLineup() {
      return sponsorLineupJson.map(sponsor => ({
        ...sponsor,
        blueskyHandle: sponsor.blueskyHandle
          ? `https://bsky.app/profile/${entry.blueskyHandle}`
          : sponsor.blueskyHandle
      }));
    },
    groupExpertCornerTopics() {
      this.expertCornerLineupUnsorted.forEach((corner) => {
        const timeSlot = corner.startTime;
        if (!this.expertCornerLineup[timeSlot]) {
          this.expertCornerLineup[timeSlot] = [];
        }
        this.expertCornerLineup[timeSlot].push(corner);
      });
    },
    openSpeakerInfoModal(speakers, id) {
      this.activeSpeakers = speakers;
      this.$refs.agenda.ariaHidden = true;
      this.$refs.speakerModal.ariaHidden = false;
      this.$refs.speakerModal.style.display = "flex";
      this.lastFocussedElementID = id;

      setTimeout(() => {
        this.$refs.speakerModal.focus();
      }, 0);
    },
    closeSpeakerInfoModal() {
      this.activeSpeakers = null;
      this.$refs.agenda.ariaHidden = false;
      this.$refs.speakerModal.ariaHidden = true;
      this.$refs.speakerModal.style.display = "none";

      for (const key in this.$refs) {
        if (
          key.startsWith("bluesky") ||
          key.startsWith("github") ||
          key.startsWith("linkedin") ||
          key.startsWith("mastodon") ||
          key.startsWith("bluesky")
        ) {
          delete this.$refs[key];
        }
      }
      document.getElementById(this.lastFocussedElementID).focus();
    },
    focusTrapModal($event) {
      let focussableElements = [];
      focussableElements.push(this.$refs.close);

      for (const key in this.$refs) {
        if (
          key.startsWith("bluesky") ||
          key.startsWith("github") ||
          key.startsWith("linkedin") ||
          key.startsWith("mastodon") ||
          key.startsWith("bluesky")
        ) {
          const element = this.$refs[key];
          if (Array.isArray(element)) {
            focussableElements.push(element[0]);
          } else {
            focussableElements.push(element);
          }
        }
      }

      const filteredFocussableElements = focussableElements.filter(
        (el) => el !== undefined
      );
      const activeElementIndex = filteredFocussableElements.indexOf(
        $event.target
      );

      if (activeElementIndex != filteredFocussableElements.length - 1) {
        if ($event.shiftKey) {
          if (activeElementIndex === 0) {
            filteredFocussableElements[
              filteredFocussableElements.length - 1
            ].focus();
          } else {
            filteredFocussableElements[activeElementIndex - 1].focus();
          }
        } else {
          filteredFocussableElements[activeElementIndex + 1].focus();
        }
      } else {
        if ($event.shiftKey) {
          filteredFocussableElements[activeElementIndex - 1].focus();
        } else {
          filteredFocussableElements[0].focus();
        }
      }
    },
    createCalendars() {

      let newStartTime = "2025-06-05T00:00:00.000";
      let newEndTime = "2025-06-06T00:00:00.000";

      let calendarStartDate = new Date(newStartTime).toISOString().replace(/-|:|\.\d+/g, '');
      let calendarEndDate = new Date(newEndTime).toISOString().replace(/-|:|\.\d+/g, '');

      let officeStartDate = new Date(newStartTime).toISOString();
      let officeEndDate = new Date(newEndTime).toISOString();

      const forbiddenCharacters = new RegExp('#', 'g')
        const removeForbiddenCharachters = (text) => {
            if (typeof text === 'string') {
              let formattedText = text.replace(/(&amp;|&)/g, " and ");
              return formattedText.replace(forbiddenCharacters, '');
            }
            return ''
        }

        const removeForbiddenCharachtersOutlook = (text) => {
          if (typeof text === 'string') {
            let formattedText = text.replace(/(?:\r\n|\r|\n)/g, "\\n");
            return formattedText.replace(forbiddenCharacters, '');
          }
          return ''
      }

      let description = "Dear friend, \n\nWe are thrilled to announce that SAP Inside Track Netherlands 2025 will be held on November 29, 2025 . This event is the highlight of the year for developers working with UI5, and we have an exciting lineup of keynotes and sessions planned for you.\n\nFor those eager to join us in person at St. Leon-Rot, Germany, it's important to note that seats for SAP Inside Track Netherlands 2025 are limited. To secure your spot, registration will be required. Stay updated by regularly visiting our conference homepage or following us on our social media channels, so you don't miss the registration start date.\n\nCan't make it to St. Leon-Rot? Don't worry, we've got you covered! SAP Inside Track Netherlands 2025 will also feature a livestream on the OpenUI5 YouTube channel, showcasing selected sessions. This way, you won't miss out on the valuable insights and knowledge shared by our esteemed speakers. Block this date in your calendar, as you won't want to miss the opportunity to learn from the best in the industry.\n\nStay tuned for more updates and exciting announcements as we get closer to SAP Inside Track Netherlands 2025. We look forward to seeing you there, either in person or virtually!\n\nPlease save the following details:\nDate: 05/06/2025 \nLocation: SAP Nederland, Amerikastraat 10, 5232 BE \'s-Hertogenbosch (The Netherlands) \n\nConference website: https://sitnl.nl"

      let descriptionGoogle = `Dear friend, <br><br>We are thrilled to announce that SAP Inside Track Netherlands 2025 will be held on November 29, 2025 . This event is the highlight of the year for developers working with UI5, and we have an exciting lineup of keynotes and sessions planned for you.<br><br>For those eager to join us in person at St. Leon-Rot, Germany, it's important to note that seats for SAP Inside Track Netherlands 2025 are limited. To secure your spot, registration will be required. Stay updated by regularly visiting our conference homepage or following us on our social media channels, so you don't miss the registration start date.<br><br>Can't make it to St. Leon-Rot? Don't worry, we've got you covered! SAP Inside Track Netherlands 2025 will also feature a livestream on the OpenUI5 YouTube channel, showcasing selected sessions. This way, you won't miss out on the valuable insights and knowledge shared by our esteemed speakers. Block this date in your calendar, as you won't want to miss the opportunity to learn from the best in the industry.<br><br>Stay tuned for more updates and exciting announcements as we get closer to SAP Inside Track Netherlands 2025. We look forward to seeing you there, either in person or virtually!<br><br>Please save the following details:<br>Date: 05/06/2025<br>Location: SAP Nederland, Amerikastraat 10, 5232 BE \'s-Hertogenbosch (The Netherlands) <br><br>Conference website: <a href="https://openui5.org/ui5con/germany2025/" target="_blank">https://sitnl.nl</a>`


      let cal = [
        'BEGIN:VCALENDAR',
        'VERSION:2.0',
        'BEGIN:VEVENT',
        'DTSTART:' + calendarStartDate,
        'DTEND:' + calendarEndDate,
        'SUMMARY:' + 'Save the Date: SAP Inside Track Netherlands 2025',
        'LOCATION:' + 'SAP Nederland, Amerikastraat 10, 5232 BE \'s-Hertogenbosch (The Netherlands)',
        'DESCRIPTION:' + removeForbiddenCharachtersOutlook(description),
        'UID:' + '1',
        'END:VEVENT',
        'END:VCALENDAR'
      ].join('\n');

      return {
        calendars: [
          {
            google: encodeURI([
              'https://www.google.com/calendar/render',
              '?action=TEMPLATE',
              '&text=' + 'Save the Date: SAP Inside Track Netherlands 2025',
              '&dates=' + calendarStartDate ,
              '/' + calendarEndDate,
              '&location='+'SAP Nederland, Amerikastraat 10, 5232 BE \'s-Hertogenbosch (The Netherlands)',
              '&details=' + removeForbiddenCharachters(descriptionGoogle),
              '&sprop=&sprop=name:'
            ].join('')),
            office365: encodeURI([
              'https://outlook.office365.com/owa/',
              '?path=/calendar/action/compose',
              '&rru=addevent',
              '&subject=' + 'Save the Date: SAP Inside Track Netherlands 2025',
              '&startdt=' + officeStartDate,
              '&enddt=' + officeEndDate,
              '&location=' + 'SAP Nederland, Amerikastraat 10, 5232 BE \'s-Hertogenbosch (The Netherlands)',
              '&body=' + removeForbiddenCharachters(descriptionGoogle)
            ].join('')),
            ics: encodeURI('data:text/calendar;charset=utf8,' + cal)
          }
        ]
      }

    },
    formatBlueskyLink(handle) {
      if (!handle.startsWith('https:')) {
        return "https://bsky.app/profile/" + handle;
      }
    },
    formatLinkedInLink(handle) {
      if (!handle.startsWith('https:')) {
        return "https://www.linkedin.com/in/" + handle;
      }
    },
    formatMastodonLink(handle) {
      if(!handle.startsWith('https:')) {
        if (handle.includes('@saptodon.org')) {
          return 'https://saptodon.org/' + handle.replace('@saptodon.org', '');
        }

        return 'https://saptodon.org/' + handle;
      }
    },
    formatBlueskyLink(handle) {
      if(!handle.startsWith('https:')) {
        return 'https://bsky.app/profile/' + handle.replace('@', '');
      }
    },
    shuffleSpeakersArray(array) {
      const newArray = [...array]
      const filteredArray = newArray.filter((el) => el.hasPhoto);
      const length = filteredArray.length

      for (let start = 0; start < length; start++) {
        const randomPosition = Math.floor((filteredArray.length - start) * Math.random());
        const randomItem = filteredArray.splice(randomPosition, 1);
        filteredArray.push(...randomItem);
      }

      return filteredArray;
    },
    formatSpeakersArray(array) {
      const newArray = [...array]
      const formattedArray = newArray.map(speaker => {
        const fullName = speaker.firstName + ' ' + speaker.lastName;

        if(speaker.blueskyHandle) {
          speaker.blueskyHandle = this.formatBlueskyLink(speaker.blueskyHandle);
        }

        if(speaker.linkedInUrl) {
          speaker.linkedInUrl = this.formatLinkedInLink(speaker.linkedInUrl);
        }

        if(speaker.mastodonHandle) {
          speaker.mastodonHandle = this.formatMastodonLink(speaker.mastodonHandle);
        }

        if(speaker.blueskyHandle) {
          speaker.blueskyHandle = this.formatBlueskyLink(speaker.blueskyHandle);
        }

        return {
          ...speaker,
          fullName: fullName,
          showMore: false
        }
      });

      return formattedArray;
    },
    formatAndShuffleSpeakersArray(array) {
      const formattedArray = this.formatSpeakersArray(array);
      return this.shuffleSpeakersArray(formattedArray);
    },
    updateLiveSession() {
      return this.formattedLineup.map((session) => {
        let timeNow = new Date().toISOString();
        let sessionTimeStart = new Date(session.startTime).toISOString();
        let sessionTimeEnd = new Date(session.endTime).toISOString();

        if (timeNow >= sessionTimeStart && timeNow < sessionTimeEnd) {
          session.isLive = true;
        } else {
          session.isLive = false;
        }
      });
    },
  },
  filters: {
    trimTime: function (value) {
      let time = value.substring(value.indexOf("T") + 1);
      let timeSplit = time.split(":");
      let hour = timeSplit[0].startsWith("0")
        ? timeSplit[0].replace(/^0+/, "")
        : timeSplit[0];
      return hour + ":" + timeSplit[1];
    },
    convertTime: function (value, eventTime) {
      if (eventTime === "local") {
        return luxon.DateTime.fromISO(value)
          .toLocal()
          .toISO({ suppressMilliseconds: true });
      }
      return value;
    },
    formatLocation: function (value) {
      if (value) {
        if (value.includes("audimax")) {
          return "Audimax";
        } else if (value.includes("w1") || value.includes("w2")) {
          return "Room W1/W2";
        } else if (value.includes("w3")) {
          return "Room W3"
        } else if (value.includes("expert")) {
          return "Expert Corner"
        } else if (value.includes("canteen")) {
          return "Canteen"
        } else {
          return value;
        }
      }
    },
    formatType: function (value) {
      if (value) {
        if (value.includes("presentation")) {
          return "Talk";
        } else if (value.includes("hands")) {
          return "Workshop";
        } else if (value.includes("expert")) {
          return "Expert Corner";
        } else if (value.includes("keynote")) {
          return "Keynote";
        }  else {
          return "";
        }
      }
    },
  },
});
