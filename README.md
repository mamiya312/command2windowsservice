# command2windowsservice

command2windowsservice is a Windows Service application that starts an application specified by arguments.
Applications that are not implemented as Windows services can also be started as Windows services.

```
sc create notepad binPath= "C:\path\to\command2windowsservice.exe --name notepad C:\Windows\System32\notepad.exe"
```
