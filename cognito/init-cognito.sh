#!/bin/bash

# Cognitoのデータベースファイルを初期化するスクリプト

COGNITO_DIR="$(dirname "$0")/.cognito/db"

# ディレクトリが存在しない場合は作成
mkdir -p "$COGNITO_DIR"

# テンプレートからファイルをコピー
if [ -f "$COGNITO_DIR/local_pool_id.json.template" ]; then
    cp "$COGNITO_DIR/local_pool_id.json.template" "$COGNITO_DIR/local_pool_id.json"
    echo "Cognito local database initialized from template"
else
    # テンプレートがない場合は新規作成
    cat > "$COGNITO_DIR/local_pool_id.json" << 'EOF'
{
  "Id": "local_pool_id",
  "Name": "Local User Pool",
  "Policies": {
    "PasswordPolicy": {
      "MinimumLength": 8,
      "RequireUppercase": false,
      "RequireLowercase": false,
      "RequireNumbers": false,
      "RequireSymbols": false
    }
  },
  "UsernameAttributes": [
    "email"
  ],
  "MfaConfiguration": "OFF",
  "EstimatedNumberOfUsers": 0,
  "EmailVerificationMessage": "Your verification code is {####}",
  "EmailVerificationSubject": "Your verification code",
  "AutoVerifiedAttributes": [
    "email"
  ],
  "VerificationMessageTemplate": {
    "DefaultEmailOption": "CONFIRM_WITH_CODE",
    "EmailMessage": "Your verification code is {####}",
    "EmailSubject": "Your verification code"
  },
  "Users": {}
}
EOF
    echo "Cognito local database initialized"
fi