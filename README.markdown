Package tz provides a wrapper for `time.Location` with additional information
about timezones.

Every timezone is represented by a zone name (`Europe/Amsterdam`) and a country
(`NL`). The rationale for this is mostly a UI-one: it makes more sense to
display a list of countries first and then a list of timezones, instead of
just presenting a huge list of timezones.

For example, my current TZ is Asia/Makkasar (WITA), but it's much easier to
select "Indonesia" from a list, and then choose one of the 4 timezones in
Indonesia.

We need to store both to make sure people who fill in "Isle of Man,
Europe/London" aren't shown "you selected Britain, Europe/London" when they
revisit a settings page.

---

The data is generated from the [IANA time zone database 2019c][zoneinfo]; use
`go generate` to re-create it from `/usr/share/zoneinfo`.

**Caveat**: There is no way to store the full contents of `time.Location`
compile-time, so mismatches can occur if the generated data is from a different
version than than what Go uses.

The logic to load the tzdata isn't exported from the time package, and don't
really feel like copying and adapting it all, so ... do this for now anyway.
Just make sure you use the same TZ database version.

[zoneinfo]: http://www.iana.org/time-zones
