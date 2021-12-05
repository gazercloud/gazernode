# GazerNode
This is a small and simple application that runs as a Windows service to record metrics several times per second. Metrics can be very different. For example, memory usage by a process or ping to a host. The application does not require a DBMS. The data is stored in an open binary format. Data viewing is possible in the form of graphs of the history of changes and in tables of current values. The configuration setting is done without editing the config files - everything is available directly in the UI application.


# Unit types

## File system sensors
- File Content
- File Size

## General sensors
- CGI
- CGI Key=Value
- HHGTTG

## Network sensors
- Ping
- TCP Connect

## Serial Port sensors
- Serial Port Key=Value

## Windows sensors
- Windows Memory
- Windows Process
- Windows Networks
