# Smart Home Project
Smarthome thesis project

[Visualized description](https://miro.com/welcomeonboard/cTJxMlpKNDVTWW4wbndPR0dxcVJqdVFaSXZidWRiSVo0cWNaRVdpcTNOR0xlTTNVZ1NjVVpoUmloVmRnbGdKeXwzNDU4NzY0NTQ3MjMxMDAyNDg1fDI=?share_link_id=25637035199)

## Documentations
- [Developer documentation(ENG)](docs/dev_documentation_en.md)
- User documentation(ENG)
- Developer documentation(HUN)
- User documentation(HUN)

## Description
This a remote lamp controller project using [ESP32_Relay_X8](https://templates.blakadder.com/ESP32_Relay_X8.html).  
Client connect to a publicly reachable server and that communicate with the controller via websocket and provide instructions.
Due to the server is public the user have to authenticate with a password, this is managed with haproxy.
