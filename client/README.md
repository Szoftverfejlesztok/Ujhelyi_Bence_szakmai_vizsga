# Controller setup

## Config
Copy `include/config.def.h` to `include/config.h` and edit it

## Build && Upload
```
cd SmartHome
platformio run --target upload --upload-port /dev/ttyUSB0
```

## Monitor
Optional but could be handy in development
```
platformio device monitor -b 115200
```


