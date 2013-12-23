# git-largefile

Manage binary files with git.

## How it works

コミットするときにハッシュ値だけをコミットし、ファイルの実態は `~/.gitasset/data`
に格納します。

別のマシンでチェックアウトする場合は、 `~/.gitasset/data` を rsync しておきます.

## Setup

### Install

``store-largefile.py`` と ``load-largefile.py`` にパスを通してください。
`pip install path.py` もしておいてください。

### S3 Configuration

予め S3 にアクセスできるキーとバケットを作っておいてください。

`~/.gitasset/gits3.ini` に次のように書いてください:

```
[DEFAULT]
awskey = "Access Key Id:Secret Access Key"
bucket = バケット名
```

### gitconfig

`~/.gitconfig` か `.git/config` に、次のように設定してください

```
[filter "s3"]
    clean = gits3 store
    smudge = gits3 load
    required
```

### gitattribute

git リポジトリの中に `.gitattributes` っていうファイルを作って、次のように設定してください。

```
*.png  filter=s3
*.jpeg filter=s3
*.jpg  filter=s3
*.gif  filter=s3
```

これで設定したファイルは largefile フィルターを通るようになります.
