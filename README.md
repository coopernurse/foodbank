

## Run locally

You'll need:

* Go 1.23 or later
* Docker
* `gcloud` CLI

### Run Firestore emulator and set env var

Make sure `gcloud` is installed, then run:

```
gcloud emulators firestore start
```

Then set FIRESTORE_EMULATOR_HOST based on the output of the above command.  For example:

```
export FIRESTORE_EMULATOR_HOST=[::1]:8548
```

### Start server

```
make start
```

