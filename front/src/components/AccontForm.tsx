import { Dispatch, SetStateAction } from 'react';
import { Account } from "../model/account"

export type AccountFormProps = {
  setAccount: Dispatch<SetStateAction<Account>>
};
export function AccountForm({ setAccount }) {
  return <div>The form</div>
}
