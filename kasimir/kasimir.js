/* Kasimir helpers: lazy play/pause for every <video data-kasimir>, poster only for reduced motion,
   plus an optional "spot Kasimir" peek. No dependencies. */
(function () {
  var reduce = window.matchMedia && window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  var vids = document.querySelectorAll('video[data-kasimir]');

  vids.forEach(function (v) {
    v.muted = true; v.loop = true; v.playsInline = true; v.preload = 'none';
    if (reduce) { v.removeAttribute('autoplay'); return; }
  });

  if (!reduce && 'IntersectionObserver' in window) {
    var io = new IntersectionObserver(function (entries) {
      entries.forEach(function (e) {
        var v = e.target;
        if (e.isIntersecting) { if (v.preload === 'none') { v.preload = 'auto'; v.load(); } v.play().catch(function () {}); }
        else { v.pause(); }
      });
    }, { threshold: 0.25 });
    vids.forEach(function (v) { io.observe(v); });
  }

  /* Peek: add <button class="kasimir-peek" data-src="kasimir/kasimir-beard.mp4" data-delay="20000"></button> */
  var peek = document.querySelector('.kasimir-peek');
  if (!peek) return;
  var KEY = 'sitnl-kasimir-spotted';
  var spotted = 0;
  try { spotted = parseInt(localStorage.getItem(KEY) || '0', 10); } catch (e) {}
  var v = document.createElement('video');
  v.src = peek.dataset.src || 'kasimir/kasimir-idle.mp4';
  v.muted = true; v.loop = true; v.playsInline = true; v.setAttribute('aria-hidden', 'true');
  var label = document.createElement('span'); label.className = 'kasimir-peek__label';
  peek.setAttribute('aria-label', 'Kasimir spotted! Click to say hi');
  peek.append(v, label);
  peek.tabIndex = -1; /* not in the tab order while hidden off-screen */
  setTimeout(function () {
    peek.classList.add('is-visible'); peek.tabIndex = 0;
    if (!reduce) v.play().catch(function () {});
    setTimeout(function () { peek.classList.remove('is-visible'); peek.tabIndex = -1; v.pause(); }, 6000);
  }, parseInt(peek.dataset.delay || '15000', 10));
  peek.addEventListener('click', function () {
    spotted += 1;
    try { localStorage.setItem(KEY, String(spotted)); } catch (e) {}
    label.textContent = 'Spotted ' + spotted + 'x';
  });
})();
