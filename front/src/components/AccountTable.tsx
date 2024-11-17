import { Accounts, Account } from "../model/account";

type AccountTableProps = {
  accounts: Accounts;
  message: string;
};

export function AccountTable({ accounts, message }: AccountTableProps) {
  return <div>
    <table>
      <thead>
        <tr>
          <th>Number</th>
          <th>Name</th>
          <th>Iban</th>
          <th>Address</th>
          <th>Amount</th>
          <th>Type</th>
        </tr>
      </thead>
      <tbody>
        {
          message === "" ? accounts.map((acc: Account) => {
            return (<tr key={acc.number}>
              <th>{acc.number}</th>
              <th>{acc.name}</th>
              <th>{acc.iban}</th>
              <th>{acc.address}</th>
              <th>{acc.amount}</th>
              <th>{acc.type}</th>
            </tr>);
          }) : (<tr><td colSpan={6}>{message}</td></tr>)
        }
      </tbody>
    </table>
  </div>
}

