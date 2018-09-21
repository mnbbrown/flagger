
const getURL = (path: string) => {
  return process.env.REACT_APP_API ? `${process.env.REACT_APP_API}${path}` : path;
}

const parseResponse = (response: Response) : Promise<any> => {
  if (response.ok) {
    return response.json();
  }
  throw new Error("Bad response");
}

export default {
  getFlags: () : Promise<any> => {
    return fetch(getURL('/flags')).then(parseResponse);
  },
  setFlag: (flag: string, env: string, type: string, value: number) : Promise<any> => {
    return fetch(getURL(`/flags/${flag}/${env}`), {
      body: JSON.stringify({ type, value: `${value}` }),
      headers: {
        "Content-Type":"application/json"
      },
      method: "POST"
    })
  }
}
