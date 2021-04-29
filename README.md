<div align="center">
  <h1>go-richelieu</h1>
  <p>Data generator that respects cardinality and schema structure.</p>
  <img source=""/>
  <a href="https://coveralls.io/github/estebgonza/go-richelieu" title="Data generator that respects cardinality and schema structure.">
    <img src="https://github.com/estebgonza/go-richelieu/workflows/Go/badge.svg" alt="Data generator that respects cardinality and schema structure."/>
  </a>
</div>

### Usage
- Build with `make`
- Initiate a new plan.json `./bin/richelieu_linux_amd64 readFromColumn "INT, STRING, INT"`
- Adapt the plan.json file describing dataset to generate. See plan.json.example.
- Generate the dataset `./bin/richelieu_linux_amd64 generate`

### Dataset configuration in plan.json
- **Distinct:** cardinality of column values
- **Name:** column name
- **Mode:**
  - **Block:** 0, 0, 0, 1, 1, 1, 2, 2, 2
  - **Alternate:** 0, 1 , 2, 0, 1 , 2, 0, 1 , 2
  - **Random:** 2, 0, 1, 2, 2, 0, 1, 0, 1
- **Prefix:** string column prefixed with this string
- **Offset:** start index from this offset
- **Values:** force a list of values to appear in generated data
- **Start/End:** used to control interval for float and dateTime columns
- **blockstep** for block mode, the same value appears consecutively this number of times
  - **blockstep: 2** 0, 0, 1, 1, 2, 2, 0, 0, 1, 1, 2, 2, ...
  - **blockstep: 3** 0, 0, 0, 1, 1, 1, 2, 2, 2, 0, 0, 0, ...
