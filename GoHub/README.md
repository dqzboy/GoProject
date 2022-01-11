## 简易版上传文件至阿里云OSS的小工具
### 1、首先在自己项目中创建etc目录并创建一个`.secret.env`文件
### 2、将阿里云账号的AccessKey以及OSS ENDPOINT添加进去
```
export ALI_AK="AccessKey ID"
export ALI_SK="AccessKey Secret"
export ALI_OSS_ENDPOINT="Endpoint"
```
### 3、修改main.go代码中的BucketName为自己的Bucket
### 4、build之前先执行`source etc/.secret.env`
![gohub01](https://user-images.githubusercontent.com/42825450/148933882-ebf2acc1-0a8e-4cd6-afd0-d4b0daf764c9.jpeg)
![gohub02](https://user-images.githubusercontent.com/42825450/148933889-2ad1fcb3-cd1d-4c0b-91f1-60fb6d1104ef.jpeg)
