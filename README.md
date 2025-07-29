# HackTheBox Academy to Markdown
This is a simple CLI application that will fetch and convert a HackTheBox Academy module into a local file in Markdown format.
This program will only grab one module at a time, and requires authenticating with the platform. 
You will also need to have the module unlocked, which should go without saying.

I personally use [Obsidian](https://obsidian.md/) as my note-taking tool, and this application is tailored and tested for rendering markdown utilizing it.
Most other note-taking tools that can import markdown files should work fine as well.

### Disclaimer
**Please note that this application is not intended for use in uploading or sharing the end result content.**
The application is solely designed for personal use and any content created using this application should not be shared or uploaded to any platform without proper authorization and consent from HackTheBox.
The contributors of this application are not responsible for any unauthorized use or distribution of the content created using this application.

### Installing
Check the releases folder [here](https://github.com/Tut-k0/htb-academy-to-md/releases), and download the most recent executable for your operating system.
All the executables listed here are for x64 and amd64. If there is not an executable for your OS or architecture, you can simply build the application. (See building section below.)

### Running
These steps have changed slightly with the reCaptcha update on the HackTheBox platform. 
I see this current state as a workaround for not dealing with the reCaptcha until I have more time to dig into that.

Essentially instead of passing your email and password, you will just pass your authenticated session cookies to the application to use. 
So the one added step for the workaround is manually logging into the academy (I would assume you are logged into Academy anyway to get the module URL), and extracting your cookies from your browser.
You can fetch these with the developer tools, burp, or a browser extension, whatever works easiest for you.
The cookies will get passed to the new `-c` argument, and you no longer need to pass an email or password.
```bash
# Get the help menu displayed
> htb-academy-to-md -h

# Run with a file containing multiple modules, write markdown to /tmp/module and images to /tmp/module/images
> htb-academy-to-md -m /tmp/modulelist.txt -o /tmp/module/ -i /tmp/module/images/ -c "htb_academy_session="
Authenticating with HackTheBox...
Downloading module https://academy.hackthebox.com/module/136/section/1259
Downloading module images...
Finished downloading module!
Downloading module https://academy.hackthebox.com/module/116/section/1140
Downloading module images...
Finished downloading module!
Downloading module https://academy.hackthebox.com/module/115/section/1101
Downloading module images...
Finished downloading module!
```

### Building
```bash
# Run from inside the /src folder.
go build -o htb-academy-to-md
```
