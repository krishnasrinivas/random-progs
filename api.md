## s3Client creation
```js
var s3Client = new Minio({
  accessKey: 'YOUR-ACCESSKEYID',
  secretKey: 'YOUR-SECRETACCESSKEY',
  bucket: 'BUCKETNAME' // mandatory for all calls except listBuckets
  cloud: 'CLOUDNAME', // can be 'amazon', 'google' or 'minio'
  minioEndpoint: 'https://minioServer.com', // needed only if cloud is minio, 
})
```

## API signatures

### Bucket ops
listBuckets(cb)

makeBucket(region, cb) //region is optional

bucketExists(cb)

removeBucket(cb)

getBucketACL(cb)

setBucketACL(acl, cb)

listObjects(prefix, recursive)

listIncompleteUploads(prefix, recursive)

### Object ops
getObject(objectName, cb)

getPartialObject(objectName, offset, length, callback)

putObject(objectName, contentType, size, stream, callback)

statObject(objectName, callback)

removeObject(objectName, callback)

removeIncompleteUpload(objectName, callback)

### Presigned ops
presignedGetObject(objectName, expiry)

presignedPutObject(objectName, expiry)

presignedPostPolicy -> here, policy.setBucket("bucketname") is not needed

