## /api/setOwner

Update owner info of a branch.

### Parameter

- `repo`: Repository name in `namespace/project` format.
- `branch`: Branch name.
- `owner`: Owner info.

### Response Status

- `200`: Success.
- `404`: No such branch or project.
- `500`: Server error like lost db connection.

## /api/setDesc

Update description of a branch.

### Parameter

- `repo`: Repository name in `namespace/project` format.
- `branch`: Branch name.
- `desc`: Description.

### Response Status

- `200`: Success.
- `404`: No such branch or project.
- `500`: Server error like lost db connection.

## /api/list

Retrieve all project info

### Response Status

- `200`: Success.
- `500`: Server error like lost db connection.

### Response Body

```json
[{
  "name": "namespace/project",
  "branches": [{
    "name": "branch_name",
    "owner": "owner info",
    "desc": "description"
  }]
}]
```