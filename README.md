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
Testing is pretty straight forward:
  1. Set up a calendar event on Google Calendar. Make sure it's for the account
     listed in settings.
  2. Run the application with `/bin/xantoria.com`
  3. A notification should appear at the correct time. An INIT line should also
     be visible in the logs, which is a good indication things are working.
