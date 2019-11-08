
# KLF Group Junior Assessment Submission

## By Geoffrey Guest

* [redeyedmars.com](http://www.redeyedmars.com/) - The website used to demonstrate the assessment answers
* [redeyedmars.com/report.html](http://www.redeyedmars.com/report.html) - The website used to demonstrate the assessment 3

## Overview

### Techniques

I setup a web server that could host html files. This html page hosts a websocket to anyone that attempts to connect at / or /report.html

The messaging system sends packets of data in this format: {actionname:username}actiondata_0,actiondata1...

I used regexes on both sides (front and server) to parse these messages

When an action is performed a database entry is made with the time and activity type

* Note that the activity is determined based on a HashMap in memory, but that if the activity name does not have an id, one is dynamically created by adding the entry to the database and then retrieving the id of that activity.

* A valid username must be somewhere in the message for the activity to be logged, as such no logging is done on page load

A note about how the passwords are passed, this is the sequence:

-> Submission of form -> sha256 is called on the password salted with the username

-> Password is sent to the server -> Server salts -> Server hashes again

-> The password is used in its hashed form to both lookup users and is saved in that hashed form


### Technologies

* Server-side - golang
* Front-end - html/css/javascript
* encryption - [sha256 using this javascript implementation](https://github.com/brillout/forge-sha256/)
* Database - mysql
* Cloud - aws
* Version Control - github.com


### Getting started

If you would like to host your own web server using the golang code provided, there are tinkering with the code as the addresses to the webserver are hardcoded, as well as the address of the database used. The required files should be found in src/networking/web.go and src/databasing/setup.go respectively.

The server is run using the following command: go run user_server.go

## Assessment 1

API's for adding users can be seen at the front page of [redeyedmars.com](http://www.redeyedmars.com/)

Once on the site the user can change their username/password by clicking on their username in the top left.

As well as on the server's command prompt, the admin can issue the following commands:

add [username] [password]

remove [username]

change [username] [new_username] [new_password]

All changes are recorded in the mysql backend.

## Assessment 2

The main page of [redeyedmars.com](http://www.redeyedmars.com/) showcases a user login/signup page that then transitions to a barebones website

On the left is the navigation bar.

On the main page the user can view some business details.

The user can also navigate to the "Contact" page, located on the left hand side of the screen.

* Note that I could not figure out a way to hook up an email to this webserver, so I am currently just logging the contact me's into the mysql database.

The user can also logout in the bottom left of the page.

## Assessment 3

The side page of [redeyedmars.com/report.html](http://www.redeyedmars.com/report.html) can be viewed. Since data is being used from active users, and we are currently in the month of November no data is reported for October. However there is a drop down of the months that the user can select November in and the report will populate for November.

## Bonus Challenge

The Bonus Challenge was completed using the webserver + website to fully showcase the 3 assessments.
