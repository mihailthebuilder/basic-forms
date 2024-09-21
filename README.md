# basic-forms

Go webhook that stores form submissions in the filesystem.

## How it works

- Request a user, and receive an external and internal ID.
- Use the external ID to handle form submissions in your client application.
- Fetch submissions with the internal ID (keep it secret!).

## TODO

- send 404s when submission not found
- figure out why docker removes all state
