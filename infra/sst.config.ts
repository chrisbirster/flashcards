/// <reference path="./.sst/platform/config.d.ts" />

export default $config({
  app(input) {
    const cloudflareApiToken = (process.env.CLOUDFLARE_API_TOKEN || "").trim();
    return {
      name: "vutadex-infra",
      home: "aws",
      removal: input?.stage === "production" ? "retain" : "remove",
      providers: {
        aws: {
          region: (process.env.AWS_REGION ?? "us-east-1") as any,
        },
        ...(cloudflareApiToken !== ""
          ? {
              cloudflare: {
                apiToken: cloudflareApiToken,
              },
            }
          : {}),
      },
    };
  },
  async run() {
    const authHeaderName = (process.env.VUTADEX_EMAIL_API_AUTH_HEADER || "Authorization").trim();
    const sender = (process.env.VUTADEX_EMAIL_SENDER || "vutadex.com").trim();
    const from = (process.env.VUTADEX_EMAIL_FROM || "no-reply@vutadex.com").trim();
    const manualSesIdentity = (process.env.VUTADEX_MANUAL_SES_IDENTITY || "true").trim().toLowerCase() !== "false";
    const appOrigin = (process.env.VUTADEX_APP_ORIGIN || "https://app.vutadex.com").trim();
    const webDomain = (process.env.VUTADEX_WEB_DOMAIN || "vutadex.com").trim();
    const enableMarketingSite = (process.env.VUTADEX_ENABLE_MARKETING_SITE || "true").trim().toLowerCase() !== "false";

    const emailApiKey = new sst.Secret("EmailApiKey");
    const emailIdentity = manualSesIdentity
      ? undefined
      : new sst.aws.Email("VutadexEmail", {
          sender,
          dns: sst.cloudflare.dns(),
        });

    const emailApi = new sst.aws.Function("EmailApi", {
      handler: "functions/email.handler",
      runtime: "nodejs20.x",
      timeout: "15 seconds",
      memory: "256 MB",
      url: true,
      link: emailIdentity ? [emailIdentity, emailApiKey] : [emailApiKey],
      permissions: manualSesIdentity
        ? [
            {
              actions: ["ses:SendEmail", "ses:SendRawEmail"],
              resources: ["*"],
            },
          ]
        : undefined,
      environment: {
        EMAIL_SEND_AUTH_HEADER: authHeaderName,
        EMAIL_FROM: from,
      },
    });

    let website: sst.cloudflare.StaticSite | undefined;
    if (enableMarketingSite) {
      requireCloudflareConfig(webDomain);
      website = new sst.cloudflare.StaticSite("WebSite", {
        path: "../marketing",
        build: {
          command: "bun run build",
          output: "dist",
        },
        environment: {
          VITE_APP_ORIGIN: appOrigin,
        },
        domain: webDomain,
        errorPage: "index.html",
      });
    }

    return {
      region: process.env.AWS_REGION ?? "us-east-1",
      sender,
      from,
      manualSesIdentity,
      authHeaderName,
      emailApiBaseUrl: emailApi.url,
      marketingEnabled: enableMarketingSite,
      webDomain: website ? webDomain : undefined,
      webUrl: website?.url,
    };
  },
});

function requireCloudflareConfig(domain: string) {
  const accountId = (process.env.CLOUDFLARE_DEFAULT_ACCOUNT_ID || "").trim();
  const apiToken = (process.env.CLOUDFLARE_API_TOKEN || "").trim();
  if (accountId !== "" && apiToken !== "") {
    return;
  }
  throw new Error(
    `Cloudflare config missing for marketing deploy (${domain}). Set CLOUDFLARE_DEFAULT_ACCOUNT_ID and CLOUDFLARE_API_TOKEN.`,
  );
}
