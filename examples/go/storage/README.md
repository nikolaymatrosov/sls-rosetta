## Description

This expample showhs how to use storage trigger with Go function. The function is triggered by a new object
in the bucket. The function reads the object assuming it is an image and resizes it to 100x100 pixels, putting
the thumbnail to the same bucket into the `thumbnails` folder.

This function also uses the `libvips` to show how to provide binary dependencies to the function.
If you won't upload the library along with the function, you'll get the following error:

```
# pkg-config --cflags  -- vips
Package vips was not found in the pkg-config search path.
Perhaps you should add the directory containing `vips.pc'
to the PKG_CONFIG_PATH environment variable
No package 'vips' found
```

To solve the issue we need to download the library and upload it along with the function. It is important to provide 
the version of the library compatible with the OS where the function will be executed. To do so, we can look up the
OS version in the documentation and download the library from the [official repository](https://cloud.yandex.ru/docs/functions/lang/golang/).
Then we need to find lib in the Ubuntu repository and download:
* [libvips-dev](https://packages.ubuntu.com/jammy/amd64/libvips-dev/download);
* [libvips42](https://packages.ubuntu.com/jammy/amd64/libvips42/download).

To build the function, run the following command:

```bash
docker build --platform linux/amd64 \
    -t ycf-go:1.21.0 \
    -f ./Dockerfile .
mkdir ./build || true
docker run --rm \
    --platform linux/amd64 \
    -v "./function:/function" \
    -v "./build:/build" \
    ycf-go:1.21.0 \
    /bin/sh -c "cd function && ./build.sh"
```

In the build script, we build our function as plugin and the using `ldd` utility we find all the dependencies.
Then we copy all the dependencies to the `build/shared-libs` folder and archive it. The archive will be uploaded
to Object Storage as it will exceed the size limit for direct upload — 3.5 MB.


## Prerequisites

* [Terraform](https://www.terraform.io/downloads.html) >= 0.14.0
* [Go](https://golang.org/doc/install) >= 1.19
* [Yandex Cloud CLI](https://cloud.yandex.ru/docs/cli/quickstart)
* [curl](https://curl.se/download.html)

## Usage with Terraform deploy

To initialize Terraform, run the following command:

```bash
terraform -chdir=./tf init
```

To set the environment variables, run the following command:

```bash
export TF_VAR_cloud_id=b1g***
export TF_VAR_folder_id=b1g***
export YC_TOKEN=`yc iam create-token`
```

To deploy the infrastructure, run the following command and confirm the action typing `yes`:

```bash
terraform -chdir=./tf apply
```

To test the function you can upload an image to the bucket:

```bash
BUCKET=$(terraform -chdir=./tf output -raw bucket)

aws s3 cp --endpoint-url=https://storage.yandexcloud.net \
    ./image.jpg \
    s3://$BUCKET/uploads/image.jpg
```

Eventually, you'll see the thumbnail in the `thumbnail` folder of the bucket.


To destroy the infrastructure, run the following command:

```bash
for b in "bucket" "bucket_for_function"; do
    BUCKET=$(terraform -chdir=./tf output -raw $b)
    aws s3 rm --endpoint-url=https://storage.yandexcloud.net \
        s3://$BUCKET --recursive
done
terraform -chdir=./tf destroy --auto-approve
```