export interface AccountResponse {
	status: number
	message: string
	accounts: Accounts
}

export interface Accounts extends Array<Account> { }

export interface Account {
	number: number
	name: string
	iban: string
	address: string
	amount: number
	type: AccountType
}

export enum AccountType {
	sending = "sending",
	receiving = "receiving"
}
