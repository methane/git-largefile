# git-largefile

hg largefile みたいにおっきいファイルを git でも扱いたい

まだプロトタイプです。

## 動作

コミットするときにハッシュ値だけをコミットし、ファイルの実態は `~/.gitasset/data`
に格納します。

別のマシンでチェックアウトする場合は、 `~/.gitasset/data` を rsync しておきます.

## 設定方法

### インストール

``store-largefile.py`` と ``load-largefile.py`` にパスを通してください。
`pip install path.py` もしておいてください。

### gitconfig

`~/.gitconfig` か `.git/config` に、次のように設定してください

```
[filter "largefile"]
    clean = store-largefile.py
    smudge = load-largefile.py
```

これで largefile フィルターが動きます.

### gitattribute

git リポジトリの中に `.gitattributes` っていうファイルを作って、次のように設定してください。

```
*.png  filter=largefile
*.jpeg filter=largefile
*.jpg  filter=largefile
*.gif  filter=largefile
```

これで設定したファイルは largefile フィルターを通るようになります.
