import React, { useState } from 'react';
import { CreateAccount } from '../api';
import { AccountType } from "../model/account"

export function AccountForm() {
  const [sendMessage, setSendMessage] = useState("");
  return (
    <div>
      <form onSubmit={(e: React.FormEvent) => {
        e.preventDefault();

        const form = e.target as typeof e.target & {
          number: { value: number };
          name: { value: string };
          iban: { value: string };
          address: { value: string };
          amount: { value: number };
          type: { value: AccountType }
        }


        CreateAccount({
          // FIXME: This is pretty strange even though it is number it becomes string from form reading.
          // See: https://stackoverflow.com/questions/38703780/how-to-preserve-numeric-fields-in-json-stringify
          number: form.number.value * 1,
          name: form.name.value,
          iban: form.iban.value,
          address: form.address.value,
          amount: form.amount.value * 1,
          type: form.type.value
        }).then((response) => {
          setSendMessage("sending...")
          if (response.ok) {
            setSendMessage("")
            // FIXME: I could reload the effect from main instead to
            // just reload hooks, but works ok.
            window.location.reload()
          } else {
            setSendMessage("Failed to send account data. Try again later.")
          }
        })

      }}
      >
        <label> Unique number:
          <input type="number" defaultValue="6" name="number" />
        </label ><br />
        <label> Name:
          <input type="text" defaultValue="Your Name" name="name" />
        </label ><br />
        <label> Iban:
          <input type="text" defaultValue="PL10105097603123" name="iban" />
        </label ><br />
        <label> Address:
          <input type="text" defaultValue="123 Sesame Street" name="address" />
        </label ><br />
        <label> Amount:
          <input type="number" defaultValue="1234" name="amount" />
        </label ><br />
        <label> Type:
          <select name="type">
            <option value="sending">sending</option>
            <option value="receiving">receiving</option>
          </select >
        </label ><br />
        <button type="submit">Submit</button>
      </form >
      {sendMessage}
    </div>
  );

}
