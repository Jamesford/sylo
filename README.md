# sylo
### (Go) Sort Your Labels Out

Wasting time setting up labels on GitHub? Automate it.

Create a `labels.yml` file and run `$ sylo`

Sample `labels.yml` file
```
- name: 'Status: Backlog'
  color: 'ebfaeb'

- name: 'Status: Next'
  color: 'adebad'

- name: 'Status: In Progress'
  color: '85e085'
```

## Install
1. Download binary
2. Place binary into `$PATH`
3. `$ sylo`

## Build
1. Clone repo
2. `$ go build -o sylo main.go`
3. `$ sylo`