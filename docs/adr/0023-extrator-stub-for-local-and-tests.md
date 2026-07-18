# Extrator stub for local job and tests; Gemini behind the same port

Application and tests depend only on the Extrator port. For local `RunDailyJob` development and automated tests without `GEMINI_API_KEY`, DI may wire a stub/fake Extrator (fixture candidatos or controlled empty/error responses). Production wiring uses the Gemini adapter against the current unversioned prompt/schema. Choosing the implementation is a composition concern in `cmd` / presentation — never a branch inside domain or use-case logic.
