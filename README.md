# map of events

## API Specification

### api

`/api/`

#### v1

`/api/v1/`

##### competitor_type

GET `/api/v1/competitor_type`:

Returns all active competitors.

Response:

```json
[
  {
    "id": "dsfsfd3e-13edfsfsdf",
    "Name": "Competitor Name 1"
  },
  {
    "id": "dsfsfd3e-13edf124vdf",
    "name": "Competitor Name 2"
  }
]
```

POST `/api/v1/competitor_type`:

Creates a new competitor type

Request:

```json
"Competitor Name"
```

Response:

```json
{
  "id": "4wr432r3-52423423",
  "name": "Competitor Name"
}
```

##### founding_range

GET `/api/v1/founding_range`:

Returns a maximal available range.

Response:

```json
{
  "id": "unused",
  "low": 0,
  "high": 15
}
```

GET `/api/v1/founding_range/{id}`:

Returns a value of range with id = id

Response:

```json
{
  "id": "id",
  "low": 0,
  "high": 5
}
```

##### co_founding_range

Same as founding_range

##### organizer_level

GET `api/v1/organizer_level`:

Returns a list of all organizer levels that exists.

Response:

```json
[
  {
    "id": "adadad",
    "name": "name level",
    "code": "FDD"
  },
  {
    "id": "adadad",
    "name": "name level",
    "code": "FDD"
  },
  {
    "id": "adadad",
    "name": "name level",
    "code": "FDD"
  },
  {
    "id": "adadad",
    "name": "name level",
    "code": "FDD"
  },
  {
    "id": "adadad",
    "name": "name level",
    "code": "FDD"
  }
]
```

POST `api/v1/organizer_level`:

Creates an organizer level and returns it.

Request:

```json
{
  "name": "sfadsad",
  "code": "xxx"
}
```

Response:

```json
{
  "id": "ad12fdb34dedf",
  "name": "sfadsad",
  "code": "xxx"
}
```

##### organizer

GET `api/v1/organizer/`:

Returns all organizers that does exist.

Response:

```json
[
  {
    "id": "di3rkogfo093owefk",
    "name": "organizer X",
    "logo": "fdfgvdbgtr",
    "level": "13fbd-f3rdfx"
  },
  {
    "id": "di3rkogfo093owefk",
    "name": "organizer X",
    "logo": "fdfgvdbgtr",
    "level": "13fbd-f3rdfx"
  },
  {
    "id": "di3rkogfo093owefk",
    "name": "organizer X",
    "logo": "fdfgvdbgtr",
    "level": "13fbd-f3rdfx"
  }
]
```

POST `api/v1/organizer/`

Creates a new organizer and returns it.

Request:

```json
{
  "name": "organizer Name",
  "logo": "Organizer logo",
  "level": "fgdsfsdfd"
}
```

Response:

```json
{
  "id": "dsffafbc-gdgffd",
  "name": "organizer Name",
  "logo": "Organizer logo",
  "level": "fgdsfsdfd"
}
```

GET `api/v1/organizer/{id}`:

Returns organizer with ID = id

Response:

```json
{
  "id": "id",
  "name": "organizer Name",
  "logo": "Organizer logo",
  "level": "fgdsfsdfd"
}
```

PUT `api/v1/organizer/{id}`:

Updates organizer with ID = id

Request:

```json
{
  "name": "organizer Name",
  "logo": "Organizer logo",
  "level": "fgdsfsdfd"
}
```

DELETE `api/v1/organizer/{id}`:

Deletes organizer with given id.

Response:

```json
{
  "id": "dsffafbc-gdgffd",
  "name": "organizer Name",
  "logo": "Organizer logo",
  "level": "fgdsfsdfd"
}
```

##### event

GET `/api/v1/event`:

Returns all current

Response:

```json
[
  "234fdg0d43fd",
  "0123reidfsj43r",
  "0xvdijk2mv",
  "032edwkmd8f7"
]
```

POST `api/v1/event`:

Create event and return it.

Request:

```json
{
  "title": "",
  "organizer": "",
  "foundingType": "",
  "foundingRangeLow": 0,
  "foundingRangeHigh": 15,
  "coFoundingRangeLow": 0,
  "coFoundingRangeHigh": 15,
  "submissionDeadline": "YYYY-MM-DD",
  "considerationPeriod": "",
  "realisationPeriod": "",
  "result": "",
  "site": "",
  "document": "",
  "internalContacts": "",
  "trl": 0,
  "competitors": [
    "id0",
    "id1",
    "id2"
  ],
  "subjects": [
    "subject_1",
    "subject_2",
    "subject_3"
  ]
}
```

Response:

```json
{
  "id": "",
  "title": "",
  "organizer": "",
  "foundingType": "",
  "foundingRangeLow": 0,
  "foundingRangeHigh": 15,
  "coFoundingRangeLow": 0,
  "coFoundingRangeHigh": 15,
  "submissionDeadline": "YYYY-MM-DD",
  "considerationPeriod": "",
  "realisationPeriod": "",
  "result": "",
  "site": "",
  "document": "",
  "internalContacts": "",
  "trl": 5,
  "competitors": [
    "id0",
    "id1",
    "id2"
  ],
  "subjects": [
    "subject_1",
    "subject_2",
    "subject_3"
  ]
}
```

GET `api/v1/event/{id}`:

Returns event info with ID = id

Response:

```json
{
  "id": "",
  "title": "",
  "organizer": "",
  "foundingType": "",
  "foundingRange": {
    "low": 0,
    "high": 15
  },
  "coFoundingRange": {
    "low": 15,
    "high": 25
  },
  "submissionDeadline": "YYYY-MM-DD",
  "considerationPeriod": "",
  "realisationPeriod": "",
  "result": "",
  "site": "",
  "document": "",
  "internalContacts": "",
  "trl": 10,
  "competitors": [
    "id0",
    "id1",
    "id2"
  ],
  "subjects": [
    "subject_1",
    "subject_2",
    "subject_3"
  ]
}
```

GET `api/v1/minimal_event/`:

Returns all events in minimal version

```json
[
  {
    "id": "",
    "title": "",
    "organizer": "",
    "submissionDeadline": "",
    "trl": 0
  },
  {
    "id": "",
    "title": "",
    "organizer": "",
    "submissionDeadline": "",
    "trl": 0
  },
  {
    "id": "",
    "title": "",
    "organizer": "",
    "submissionDeadline": "",
    "trl": 0
  },
  {
    "id": "",
    "title": "",
    "organizer": "",
    "submissionDeadline": "",
    "trl": 0
  },
  {
    "id": "",
    "title": "",
    "organizer": "",
    "submissionDeadline": "",
    "trl": 0
  },
  {
    "id": "",
    "title": "",
    "organizer": "",
    "submissionDeadline": "",
    "trl": 0
  },
  {
    "id": "",
    "title": "",
    "organizer": "",
    "submissionDeadline": "",
    "trl": 0
  }
]
```

GET `api/v1/minimal_event/{id}/`:

Returns a minimal info about event with ID = id

Response:

```json
{
  "id": "",
  "title": "",
  "organizer": "",
  "submissionDeadline": "",
  "trl": 0
}
```

POST `api/v1/event/{id}`:

Updates event with ID = id

Request:

```json
{
  "title": "",
  "organizer": "",
  "foundingType": "",
  "foundingRangeLow": 0,
  "foundingRangeHigh": 15,
  "coFoundingRangeLow": 0,
  "coFoundingRangeHigh": 15,
  "submissionDeadline": "YYYY-MM-DD",
  "considerationPeriod": "",
  "realisationPeriod": "",
  "result": "",
  "site": "",
  "document": "",
  "internalContacts": "",
  "trl": 5,
  "competitors": [
    "id0",
    "id1",
    "id2"
  ],
  "subjects": [
    "subject_1",
    "subject_2",
    "subject_3"
  ]
}
```

DELETE `api/v1/event/{id}`

Deletes the event with id

PUT `api/v1/event/{id}`:

Updates info about event and returns a new value

Request:

```json
{
  "title": "",
  "organizer": "",
  "foundingType": "",
  "foundingRangeLow": 0,
  "foundingRangeHigh": 15,
  "coFoundingRangeLow": 0,
  "coFoundingRangeHigh": 15,
  "submissionDeadline": "YYYY-MM-DD",
  "considerationPeriod": "",
  "realisationPeriod": "",
  "result": "",
  "site": "",
  "document": "",
  "internalContacts": "",
  "trl": 0,
  "competitors": [
    "id0",
    "id1",
    "id2"
  ],
  "subjects": [
    "subject_1",
    "subject_2",
    "subject_3"
  ]
}
```

Response:

```json
{
  "id": "",
  "title": "",
  "organizer": "",
  "foundingType": "",
  "foundingRangeLow": 0,
  "foundingRangeHigh": 15,
  "coFoundingRangeLow": 0,
  "coFoundingRangeHigh": 15,
  "submissionDeadline": "YYYY-MM-DD",
  "considerationPeriod": "",
  "realisationPeriod": "",
  "result": "",
  "site": "",
  "document": "",
  "internalContacts": "",
  "trl": 10,
  "competitors": [
    "id0",
    "id1",
    "id2"
  ],
  "subjects": [
    "subject_1",
    "subject_2",
    "subject_3"
  ]
}
```

##### image

GET `api/v1/image/:link`:

Returns an image(base64) by that link

POST `api/v1/image/`:

Upload an image in base64, will return a link of it.

#### v2

##### organizer

GET to `apiv2/organizer`:

Get all organizers with nested images

Request:

```json
[
  {
    "id": "adadsad",
    "name": "OrganizerName",
    "logo": "data:/big-big-image",
    "level": "dbggds-2gfvdffdgv-fdfd"
  },
  {
    "id": "adadsad",
    "name": "OrganizerName",
    "logo": "data:/big-big-image",
    "level": "dbggds-2gfvdffdgv-fdfd"
  },
  {
    "id": "adadsad",
    "name": "OrganizerName",
    "logo": "data:/big-big-image",
    "level": "dbggds-2gfvdffdgv-fdfd"
  }
]
```

POST to `api/v2/organizer`:

Request:

```json
{
  "name": "OrganizerName",
  "logo": "data:/big-big-image",
  "level": "dbggds-2gfvdffdgv-fdfd"
}
```

Response:

```json
{
  "id": "adadsad",
  "name": "OrganizerName",
  "logo": "data:/big-big-image",
  "level": "dbggds-2gfvdffdgv-fdfd"
}
```

GET `api/v2/organizer/{id}`

Get organizer info with nested image

Response:

```json
{
  "id": "adadsad",
  "name": "OrganizerName",
  "logo": "data:/big-big-image",
  "level": "dbggds-2gfvdffdgv-fdfd"
}
```
