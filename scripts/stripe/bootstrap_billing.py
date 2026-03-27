#!/usr/bin/env python3

import argparse
import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path


PRODUCT_ID = "vutadex_subscription"
PRODUCT_NAME = "Vutadex Subscription"
PRODUCT_DESCRIPTION = "Workspace subscription for Vutadex"
PRICES = {
    "VUTADEX_STRIPE_BILLING_PRICE_PRO_MONTHLY": {
        "lookup_key": "vutadex_pro_monthly",
        "nickname": "Pro Monthly",
        "unit_amount": "1200",
    },
    "VUTADEX_STRIPE_BILLING_PRICE_TEAM_MONTHLY": {
        "lookup_key": "vutadex_team_monthly",
        "nickname": "Team Monthly",
        "unit_amount": "800",
    },
}


def read_env_value(env_file: Path, key: str) -> str:
    if not env_file.exists():
        return ""
    for line in env_file.read_text().splitlines():
        if line.startswith(f"{key}="):
            return line.split("=", 1)[1].strip()
    return ""


def upsert_env(env_file: Path, key: str, value: str) -> None:
    lines = env_file.read_text().splitlines() if env_file.exists() else []
    prefix = f"{key}="
    for index, line in enumerate(lines):
        if line.startswith(prefix):
            lines[index] = f"{prefix}{value}"
            break
    else:
        lines.append(f"{prefix}{value}")
    env_file.write_text("\n".join(lines) + "\n")


class StripeAPI:
    def __init__(self, secret_key: str):
        self.secret_key = secret_key

    def request(self, method: str, path: str, params: list[tuple[str, str]] | None = None) -> dict:
        url = "https://api.stripe.com" + path
        data = None
        if params:
            encoded = urllib.parse.urlencode(params, doseq=True)
            if method.upper() == "GET":
                url += ("&" if "?" in url else "?") + encoded
            else:
                data = encoded.encode("utf-8")
        req = urllib.request.Request(
            url,
            data=data,
            method=method.upper(),
            headers={
                "Authorization": f"Bearer {self.secret_key}",
                "Content-Type": "application/x-www-form-urlencoded",
            },
        )
        try:
            with urllib.request.urlopen(req) as resp:
                return json.loads(resp.read().decode("utf-8"))
        except urllib.error.HTTPError as exc:
            body = exc.read().decode("utf-8", errors="replace")
            try:
                payload = json.loads(body)
            except json.JSONDecodeError:
                payload = {"error": {"message": body}}
            message = payload.get("error", {}).get("message", f"Stripe API error {exc.code}")
            raise RuntimeError(message) from exc

    def ensure_product(self) -> dict:
        try:
            return self.request("GET", f"/v1/products/{urllib.parse.quote(PRODUCT_ID, safe='')}")
        except RuntimeError as exc:
            if "No such product" not in str(exc):
                raise
        return self.request(
            "POST",
            "/v1/products",
            [
                ("id", PRODUCT_ID),
                ("name", PRODUCT_NAME),
                ("description", PRODUCT_DESCRIPTION),
            ],
        )

    def find_price(self, lookup_key: str) -> dict | None:
        payload = self.request(
            "GET",
            "/v1/prices",
            [
                ("lookup_keys[]", lookup_key),
                ("active", "true"),
                ("limit", "1"),
            ],
        )
        data = payload.get("data") or []
        return data[0] if data else None

    def create_price(self, lookup_key: str, nickname: str, unit_amount: str) -> dict:
        return self.request(
            "POST",
            "/v1/prices",
            [
                ("product", PRODUCT_ID),
                ("currency", "usd"),
                ("unit_amount", unit_amount),
                ("nickname", nickname),
                ("lookup_key", lookup_key),
                ("recurring[interval]", "month"),
            ],
        )


def main() -> int:
    parser = argparse.ArgumentParser(description="Create or locate Stripe billing prices for Vutadex.")
    parser.add_argument("--env-file", default=".env.staging")
    parser.add_argument("--secret-key", default="")
    args = parser.parse_args()

    env_file = Path(args.env_file)
    env_file.parent.mkdir(parents=True, exist_ok=True)

    secret_key = (
        args.secret_key.strip()
        or os.environ.get("VUTADEX_STRIPE_SECRET_KEY", "").strip()
        or os.environ.get("STRIPE_SECRET_KEY", "").strip()
        or read_env_value(env_file, "VUTADEX_STRIPE_SECRET_KEY")
    )
    if not secret_key:
        print(
            "Missing Stripe secret key. Set VUTADEX_STRIPE_SECRET_KEY in the environment or env file first.",
            file=sys.stderr,
        )
        return 1

    stripe = StripeAPI(secret_key)
    product = stripe.ensure_product()

    print(f"Using Stripe product {product['id']} ({product.get('name', PRODUCT_NAME)})")
    for env_key, price_def in PRICES.items():
        price = stripe.find_price(price_def["lookup_key"])
        if price is None:
            price = stripe.create_price(
                price_def["lookup_key"],
                price_def["nickname"],
                price_def["unit_amount"],
            )
        price_id = price["id"]
        upsert_env(env_file, env_key, price_id)
        print(f"{env_key}={price_id}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
