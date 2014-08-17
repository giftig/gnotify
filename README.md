Google Calendar Project
=======================

Purpose
-------
This project aims to integrate with the Google Calendar API, synchronise events
with configured accounts regularly, and provide notification functionality via
a variety of methods including:

  * The notify-send command, to produce a nice notification pop up in linux
  * Playing a sound if appropriate
  * Sending a notification to raspbmc


Structure
---------
Google Calendar events will be stored in redis with the (configurable) key

  gcal:event:[event-type]:[event-id]

This allows some persistence so that we don't lose stuffs when the application
restarts. It will synchronise this data with the google calendar via its API
every hour (configurable).

Every minute (also configurable), we check if there's an event for which a
notification is required, and display the notification as needed. The event
will be marked as notified.
