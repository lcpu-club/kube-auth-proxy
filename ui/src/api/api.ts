import { MessagePlugin } from "tdesign-vue-next";
import { getToken } from "./token";

async function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

class Client {
  private readonly baseUrl: string;

  public username: string | null = null;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async req(
    path: string,
    opts: RequestInit,
    ignoreNullUsername = false
  ): Promise<Response> {
    if (!ignoreNullUsername) {
      this.ensureUsername();
    }

    const url = this.username
      ? (this.baseUrl + path).replace("{!NAMESPACE}", `u-${this.username}`)
      : this.baseUrl + path;

    if (opts.body && this.username && typeof opts.body === "string") {
      opts.body = opts.body.replace("{!NAMESPACE}", `u-${this.username}`);
    }

    opts.headers = {
      ...opts.headers,
      Authorization: "Bearer " + getToken(),
    };
    try {
      let resp = await fetch(url, opts);
      if (resp.status > 299) {
        throw await resp.json();
      }
      return resp;
    } catch (e: any) {
      console.error(e);
      // TODO: Better handling
      // alert("Failed to fetch " + url);
      if (e.message)
        MessagePlugin.error(e.message);
      throw e;
    }
  }

  async get(path: string): Promise<Response> {
    return await this.req(path, {
      method: "GET",
    });
  }

  async post(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: "POST",
      body: JSON.stringify(body),
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  async put(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: "PUT",
      body: JSON.stringify(body),
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  async patch(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: "PATCH",
      body: JSON.stringify(body),
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  async delete(path: string): Promise<Response> {
    return await this.req(path, {
      method: "DELETE",
    });
  }

  async deleteWithBody(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: "DELETE",
      body: JSON.stringify(body),
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  async ensureUsername() {
    if (!getToken()) window.location.href = "../oauth/redirect";
    if (this.username) return;
    
    const userInfo = await (
      await this.req(
        "/_/whoami",
        {
          method: "GET",
        },
        true
      )
    ).json();
    this.username = userInfo.username;
  }
}

export const client = new Client("../..");
