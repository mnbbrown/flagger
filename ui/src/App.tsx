import * as React from 'react';
import api from "./api";
import './App.css';
import Toggle from "./Toggle";

interface IAppState {
  flags: any;
}

interface IFlag {
  name: string;
  envs: any;
}

interface IFlagState {
  name: string
  env?: string
  value: number
  type: string
}

class App extends React.Component<any, IAppState> {

  public state: IAppState = {
    flags: {}
  };

  public componentDidMount() {
    api.getFlags().then((flags) => this.setState({ flags }))
  }

  public setPercentFlag(flag: string, value: number, env: string = "default") {
    this.setValue(flag, env, value);
    api.setFlag(flag, env, "percent", value);
  }

  public setValue(flag: string, env: string, value: number) {
      this.setState({ flags: {
        ...this.state.flags,
        [flag]: {
          ...this.state.flags[flag],
          [env]: {
            ...this.state.flags[flag][env],
            value
          }
        }
      }})
  }

  public setBoolFlag(flag: string, value: number, env :string = "default") {
    api.setFlag(flag, env, "bool", value).then(() => {
      this.setValue(flag, env, value);
    });
  }

  public render() {
    const { flags } = this.state;
    const rows = Object.keys(flags).reduce((rws : IFlag[], flag: string) => {
      return [...rws, {name: flag, envs: flags[flag]} ]
    }, []);

    return (
      <div className="App">
        <header className="App-header">
          <h1 className="App-title">flagger UI</h1>
        </header>
        <section className="App__content">
          <table>
            <thead>
              <tr>
                <th className="Flag__name">Flags</th><th>&nbsp;</th>
              </tr>
            </thead>
            <tbody>
              {this.renderRows(rows)}
            </tbody>
          </table>
        </section>
      </div>
    );
  }

  private getValue(state: IFlagState) {
    if (state.type === "BOOL") {
      const onChange = (value: number) => {
        this.setBoolFlag(state.name, value, state.env);
      }
      return <Toggle value={state.value !== 0} onChange={onChange} />
    }
    if (state.type === "PERCENT") {
      const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        this.setPercentFlag(state.name, parseInt(e.target.value, 10), state.env);
      };
      return <div><input type="range" min="0" max="100" className="Flag__input" value={state.value} onChange={onChange}/><span>{state.value}</span></div>
    }
    return <div />
  }

  private renderEnvs(flag: string, envs: any) {
    const rows = Object.keys(envs).map((envName: string) => {
      const state = this.getValue({ name: flag, ...envs[envName]});
      return (
        <tr className={`Flag--${envName === "default" ? "default" : "env"}`} key={`${flag}-${envName}`}>
          <td className="Flag__name">{envName === "default" ? flag : <span className="Env--tag">{envName}</span>}</td>
          <td className="Flag__state">{state}</td>
        </tr>
      )
    });
    return rows;
  }

  
  private renderRows(rows: IFlag[]) {
    return rows.reduce((rs: any, row: IFlag) => (
      [...rs, ...this.renderEnvs(row.name, row.envs)]
    ), []);
  }
}

export default App;