// Scheduled rebuild: every 4 hours Netlify calls this function, which triggers
// the "Scheduled data refresh" build hook. The build runs agenda/main.go
// (see netlify.toml), so each run publishes fresh event settings, agenda,
// speakers and sponsors. If the planning tool is unreachable the build fails
// and the previous deploy stays live.
//
// BUILD_HOOK_URL is set in the Netlify UI (Project configuration > Environment
// variables, scope Functions); it is not stored in git.
export default async () => {
  const hook = process.env.BUILD_HOOK_URL;
  if (!hook) {
    console.error("BUILD_HOOK_URL is not set; skipping scheduled rebuild");
    return;
  }
  const res = await fetch(hook + "?trigger_title=Scheduled+data+refresh", { method: "POST" });
  console.log(`Build hook answered ${res.status}`);
};

export const config = {
  schedule: "17 */4 * * *",
};
