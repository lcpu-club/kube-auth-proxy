import { getToken } from "./token";

async function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

class Client {
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async req(path: string, opts: RequestInit): Promise<Response> {
    const url = this.baseUrl + path;
    opts.headers = {
      ...opts.headers,
      'Authorization': 'Bearer ' + getToken()
    };
    try {
      let resp = await fetch(url, opts);
      return resp;
    } catch (e) {
      console.error(e);
      // TODO: Better handling
      alert('Failed to fetch ' + url);
      throw e;
    }
  }

  async get(path: string): Promise<Response> {
    return await this.req(path, {
      method: 'GET'
    });
  }

  async post(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: 'POST',
      body: JSON.stringify(body),
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  async put(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: 'PUT',
      body: JSON.stringify(body),
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  async patch(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: 'PATCH',
      body: JSON.stringify(body),
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }

  async delete(path: string): Promise<Response> {
    return await this.req(path, {
      method: 'DELETE'
    });
  }

  async deleteWithBody(path: string, body: any): Promise<Response> {
    return await this.req(path, {
      method: 'DELETE',
      body: JSON.stringify(body),
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }
}

export const client = new Client('../..');
