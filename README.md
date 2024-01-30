# image-service

## Routes

_*GET /upload*_  

req: form-data 
```
image: <file>
uploader: "userId"
usage: 1 //1 is for story, check protobuf for enums
```  

res: 
```json
{
  "message": "",
  "url": "url",
  "imageId": "id"
}
```
