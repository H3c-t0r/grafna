import React, { PureComponent, SyntheticEvent } from 'react';
import config from 'app/core/config';
import { UserSignup } from './UserSignup';
import { LoginServiceButtons } from './LoginServiceButtons';

const oauthEnabled = () => Object.keys(config.oauth).length > 0;
export interface Props {
  disabledLoginForm: boolean;
  passwordHint: string;
  loginHint: string;
  disableUserSignup: boolean;
  oauthEnabled: boolean;
  ldapEnabled: boolean;
  authProxyEnabled: boolean;
  onSubmit: Function;
}

export interface State {
  loggingIn: boolean;
  user: string;
  password: string;
  email: string;
  valid: boolean;
}

export class LoginForm extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      user: '',
      password: '',
      email: '',
      loggingIn: false,
      valid: false,
    };

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  handleSubmit(e: SyntheticEvent) {
    e.preventDefault();

    if (this.validate()) {
      console.log(e);
    }
  }

  handleChange(e: any) {
    this.validate();
    console.log(e.target);
    // @ts-ignore
    this.setState({
      [e.target.name]: e.target.value,
    });
  }

  validate() {
    if (this.state.user.length > 0 && this.state.password.length > 0) {
      this.setState({
        valid: true,
      });
      return true;
    } else {
      this.setState({
        valid: false,
      });
      return false;
    }
  }

  renderLoginForm() {
    return !this.props.disabledLoginForm ? (
      <form name="loginForm" className="login-form-group gf-form-group">
        <div className="login-form">
          <input
            type="text"
            name="user"
            className="gf-form-input login-form-input"
            required
            // ng-model="formModel.user"
            placeholder={config.loginHint}
            aria-label="Username input field"
            onChange={this.handleChange}
          />
        </div>
        <div className="login-form">
          <input
            type="password"
            name="password"
            className="gf-form-input login-form-input"
            required
            ng-model="formModel.password"
            id="inputPassword"
            placeholder={config.passwordHint}
            aria-label="Password input field"
            onChange={this.handleChange}
          />
        </div>
        <div className="login-button-group">
          {!this.state.loggingIn ? (
            <button
              type="submit"
              aria-label="Login button"
              className={`btn btn-large p-x-2 ${this.state.valid ? 'btn-primary' : 'btn-inverse'}`}
              onClick={this.handleSubmit}
              disabled={!this.state.valid}
            >
              Log In
            </button>
          ) : (
            <button type="submit" className="btn btn-large p-x-2 btn-inverse btn-loading">
              Logging In<span>.</span>
              <span>.</span>
              <span>.</span>
            </button>
          )}

          {config.ldapEnabled || config.authProxyEnabled ? null : (
            <div className="small login-button-forgot-password" ng-hide="ldapEnabled || authProxyEnabled">
              <a href="user/password/send-reset-email">Forgot your password?</a>
            </div>
          )}
        </div>
      </form>
    ) : null;
  }

  render() {
    return (
      <div className="login-out-box">
        <div className="login-inner-box" id="login-view">
          {config.disableLoginForm ? null : this.renderLoginForm()}
          {oauthEnabled() ? (
            <>
              <div className="text-center login-divider">
                <div>
                  <div className="login-divider-line" />
                </div>
                <div>
                  <span className="login-divider-text">{config.disableLoginForm ? null : <span>or</span>}</span>
                </div>
                <div>
                  <div className="login-divider-line" />
                </div>
              </div>
              <div className="clearfix" />

              <LoginServiceButtons />
            </>
          ) : null}

          <UserSignup />
        </div>
      </div>
    );
  }
}
