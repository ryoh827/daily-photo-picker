# daily-photo-picker

写真フォルダからランダムに1枚を選ぶCLIツール。
使用済みの写真は記録され、次回以降除外される。新しい写真が追加されると自動的に候補に含まれる。

## セットアップ

```bash
go build -o daily-photo-picker .
```

## 写真の配置

ビルドした実行ファイルと同じディレクトリに `photos/` フォルダを作成し、写真を配置する。サブフォルダにも対応。

```
daily-photo-picker
photos/
├── 2024/
│   ├── photo1.jpg
│   └── photo2.png
└── trip/
    └── photo3.heic
```

## 使い方

```bash
./daily-photo-picker          # ランダムに1枚選んでパスを表示
./daily-photo-picker status   # 状態を表示（全写真数・使用済み・残り）
./daily-photo-picker reset    # 使用済みリストをリセット
./daily-photo-picker help     # ヘルプ
```

## 仕様

- 実行するたびに `photos/` を全スキャンし、未使用の写真からランダムに1枚選択
- 使用済みの写真は `used.json` に記録
- 新しい写真を `photos/` に追加すれば自動的に候補に含まれる
- 削除された写真は使用済みリストから自動除去
- 全写真を使い切った場合は通知し、`reset` でリセット可能
- 対応形式: jpg, jpeg, png, gif, bmp, webp, heic, heif, tiff
