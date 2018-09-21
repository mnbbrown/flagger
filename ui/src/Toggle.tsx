import * as React from 'react';
import "./Toggle.css";

interface IToggleProps {
  value: boolean,
  onChange?: (value: number) => void,
}

export default class Toggle extends React.Component<IToggleProps> {

  public render () {
    return (
      <label className="Toggle">
        <input type="checkbox" checked={this.props.value || false} onChange={this.onChange} />
        <div className="Toggle__slider" />
      </label>
    );
  }

  private onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (this.props.onChange) {
      this.props.onChange(e.target.checked ? 1 : 0);
    }
  }
}