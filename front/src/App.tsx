import { useState, useEffect } from 'react';
import './App.css';
import { ListAccounts } from './api';
import { AccountTable } from './components/AccountTable';
import { AccountResponse, Account, Accounts } from './model/account';
import { AccountForm } from './components/AccontForm';

function App() {
  const [accounts, setAccounts] = useState<Accounts>(Array<Account>());
  const [loadMessage, setMessage] = useState("");
  const [inputAccount, setInputAccount] = useState<Accounts>(Array<Account>());

  const baseErrMsg = "Failed to load data, try again later."

  useEffect(() => {
    ListAccounts().then((response) => {
      setMessage("loading...")
      if (response.ok) {
        response.json().then(
          (data) => {
            const accResp = data as AccountResponse
            setAccounts(accResp.accounts)
            setMessage("")
          })
      } else {
        setMessage(baseErrMsg)
        response.json().then(
          (data: AccountResponse) => {
            setMessage(baseErrMsg + "Error message:" + data.message)
          })
      }
    }).catch(err => {
      console.log(err)
      setMessage(baseErrMsg)
    })
  }, [])

  return (
    <div className="App">
      <AccountTable accounts={accounts} message={loadMessage} />
      <AccountForm setAccount={setInputAccount} />
    </div>
  );
}

export default App;
