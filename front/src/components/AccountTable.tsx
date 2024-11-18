import { Accounts, Account } from "../model/account";
import "./AccountTable.css"

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
      <tbody >
        {
          message === "" ? accounts.map((acc: Account) => {
            return (<tr key={acc.number}>
              <td>{acc.number}</td>
              <td>{acc.name}</td>
              <td>{acc.iban}</td>
              <td>{acc.address}</td>
              <td>{acc.amount}</td>
              <td>{acc.type}</td>
            </tr>);
          }) : (<tr><td colSpan={6}>{message}</td></tr>)
        }
      </tbody>
    </table>
  </div>
}

