import { Account } from "./model/account";

const BACKEND_ADDRESS = process.env.REACT_APP_BACKEND_ADDRESS || "http://localhost:1323"

const getApiUrl = (endpoint: string) => {
  return BACKEND_ADDRESS + "/api/v1" + endpoint
}
export async function ListAccounts() {
  return await fetch(getApiUrl("/account"), {
    method: "GET",
  });
}

export async function CreateAccount(account: Account) {
  return await fetch(getApiUrl("/account"), {
    method: "POST",
    body: JSON.stringify(account)
  });
}
