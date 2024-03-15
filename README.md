# image-service

- RESTful APIs running on `localhost:80`
- gRPC APIs running on `localhost:1111`

---

### Upload images

<details>
<summary><code>POST</code> <code><b>/upload</b></code> <code>(Upload image to image server)</code></summary>

##### Body (form-data)

> | key      | required | data type | description                       |
> | -------- | -------- | --------- | --------------------------------- |
> | image    | true     | file      | The content type should be image. |
> | uploader | true     | text      | Author's MongoDB object ID        |
> | usage    | true     | text      | "1" for story                     |

##### Responses

> | http code    | content-type       | response                                                                                |
> | ------------ | ------------------ | --------------------------------------------------------------------------------------- |
> | `200`        | `application/json` | `{"message": "Success", "Url": "url of the image", "imageId": "ObjectId of the image"}` |
> | `400`, `500` | `application/json` | `{"message": "Failed", "error":"Error messages"}`                                       |

</details>
