Package tz provides a wrapper for `time.Location` with additional information
about timezones.

Every timezone is represented by a zone name (`Europe/Amsterdam`) and a country
(`NL`). The rationale for this is mostly a user interface one: it makes more
sense to display a list of countries first and then a list of timezones, instead
of just presenting a huge list of timezones.

For example, my current TZ is Asia/Makkasar (WITA), but it's much easier to
select "Indonesia" from a list, and then choose one of the four timezones in
Indonesia:

    Indonesia: Asia/Jayapura  (WIT)  – New Guinea (West Papua / Irian Jaya); Malukus/Moluccas
    Indonesia: Asia/Makassar  (WITA) – Borneo (east, south); Sulawesi/Celebes, Bali, Nusa Tengarra; Timor (west)
    Indonesia: Asia/Pontianak (WIB)  – Borneo (west, central)
    Indonesia: Asia/Jakarta   (WIB)  – Java, Sumatra

We need to store both to make sure people who fill in "Isle of Man,
Europe/London" aren't shown "you selected Britain, Europe/London" when they
revisit a settings page.

[zoneinfo]: http://www.iana.org/time-zones
