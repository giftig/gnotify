Gnotify
=======

About the project
-----------------

### Purpose
This project aims to integrate with the Google Calendar API, synchronise events
with configured accounts regularly, and provide notification functionality via
a variety of methods including:

  * The notify-send command, to produce a nice notification pop up in linux
  * Playing a sound if appropriate
  * Sending a notification to raspbmc
  * Additional options to come!

In future it will also provide additional sources of notifications and means of
notifying users about them.

### Structure
The application will synchronise with google calendar with a frequency defined
in the config (every 1h by default). It will grab notifications at that point
and pump them to the initNotifications task to schedule the notification events.

It's currently not incredibly smart about this: there's no persistence other
than that offered by syncing with the google calendar, we don't try to make up
for missing a notification which should've occurred when the application wasn't
running, etc. These features will be added in due time.

### API

#### REST

The REST notification source allows notifications to be directly passed to the
application by hitting a REST endpoint. There are two endpoints associated with
triggering notifications:

  * POSTing to /notify/trigger/ tells the application to store a new
    notification which has just been created by the client.
  * POSTing to /notify/route/ is designed to be internal to the application and
    represents a notification triggered elsewhere and passed to this node

The main difference between these endpoints is that the former marks the
notification as being a newly-triggered event, which is of type "API", and it
may be another application hitting the endpoint. The latter indicates that another
node is passing on the event, and that it either believes the receiver to be the
correct destination, or else it considers the receiver to be an authority on
how to deliver it to the correct destination.

Notifications should be `POST` requests with content type `application/json`
matching the following format:

    {
      "id": "unique_amongst_rest_notifications",
      "title": "Remember to feed Hodor",
      "message: "Hodor needs lots of food; he's a big lad",
      "priority": 10,
      "recipient": "hodor_456",
      "complete": false,
      "time": "2000-01-01T00:00:00Z"
    }

To check for recent notifications, a non-master node will make a GET request to
the `/notify/fetch/` endpoint to retrieve a list of recent notifications.
This will return an array of notifications matching the format described above,
or will redirect to the correct master if the wrong node is asked. The
destination for which to query is described by the "dest" query parameter.

---

Development & Deployment
------------------------

### Dependencies
Go dependencies are stated in requirements.txt and installed by dependencies.sh.
Other dependencies are optional, depending on how the notification system is
configured:
  * notify-send: One method of displaying a notification
  * Other methods to come, which bring dependencies

### Testing

#### Calendar
Pretty straightforward:
  1. Set up a calendar event on Google Calendar. Make sure it's for the account
     listed in settings.
  2. Run the application with `/bin/xantoria.com`
  3. A notification should appear at the correct time. An INIT line should also
     be visible in the logs, which is a good indication things are working.

#### REST
There's a convenience script for testing this: `scripts/test/notify`. Give it
a title, message, and optionally your own ID, and it'll POST a notification
to be triggered within the next few seconds.
