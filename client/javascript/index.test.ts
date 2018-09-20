import Client from "./index";
import nock from "nock";

test("returns a value", async () => {
  nock("http://test-server").get("/flags/testFlag/testEnv").reply(200, "on").get("/flags/stringFlag/testEnv").reply(200, "enabled");
  const client = new Client("http://test-server");
  let value = await client.get("testFlag", "testEnv");
  expect(value).toBe(true);
  value = await client.get("stringFlag", "testEnv")
  expect(value).toBe("enabled");
});

test("that not found returns the default value", async () => {
  nock("http://test-server").get("/flags/testFlag/testEnv").reply(404, "on");
  const client = new Client("http://test-server");
  let value = await client.get("testFlag", "testEnv");
  expect(value).toBe(true);
})

test("connection refused", async () => {
  nock("http://test-server").get("/flags/testFlag/testEnv").replyWithError({ code: "ETIMEDOUT"});
  let client = new Client("http://test-server");
  let value = await client.get("testFlag", "testEnv");
  expect(value).toBe(true);
  
  client = new Client("http://test-server", false);
  value = await client.get("testFlag", "testEnv");
  expect(value).toBe(false);
});