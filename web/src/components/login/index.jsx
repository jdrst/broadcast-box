import React, { useEffect, useState, useContext } from 'react';
import { useNavigate } from 'react-router-dom'

export const AuthContext = React.createContext(null);

export function AuthContextProvider({ children }) {
  const [authState, setAuthState] = useState({username: '', authenticated: false, pending: true})

    useEffect(() => {
      fetch('user/info')
        .then(resp => {
          if (resp.ok) {
            resp.text().then(username => {
              setAuthState({username, authenticated: true, pending: false})
            })
          } else {
            setAuthState({username: '', authenticated: false, pending: false})
          }
        })
    }, [])
    
    const logout = () => {
      return fetch('auth/logout', {method: 'POST'})
    }
    
    const login = (username, password) => {
      return fetch('auth/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        body: new URLSearchParams({username, password})
      }).then(resp => {
          if (resp.ok) {
            resp.text().then(user => {
              setAuthState({username: user, authenticated: true, pending: false})
            })
          }
          return resp
        })
    }
    
  return (
    <AuthContext.Provider value={[authState, login, logout]}>
      {children}
    </AuthContext.Provider>
  );
}

function LoginPage() {
  const navigate = useNavigate()
  const [_authState, login, _logout] = useContext(AuthContext);
  
  const postLogin = (formData) => {
    const redirectUrl = new URL(document.location).searchParams.get('redirectUrl') ?? '/';
    login(formData.get("username"), formData.get("password")).then(resp => {
      if (resp.ok) {
        navigate(redirectUrl)
      }
    })
  };
  return (
      <form className='rounded-md bg-gray-800 shadow-md p-8' method='post' action={postLogin}>
        <div className='my-4'>
          <label htmlFor='username' className='block text-sm font-bold mb-2'>Username</label>
          <input className='appearance-none border w-full py-2 px-3 leading-tight focus:outline-hidden focus:shadow-outline bg-gray-700 border-gray-700 text-white rounded-sm shadow-md placeholder-gray-200' name='username' />
        </div>
        <div className='my-4'>
          <label htmlFor='password' className='block text-sm font-bold mb-2'>Password</label>
          <input className='appearance-none border w-full py-2 px-3 leading-tight focus:outline-hidden focus:shadow-outline bg-gray-700 border-gray-700 text-white rounded-sm shadow-md placeholder-gray-200' name='password' type='password' />
        </div>
        <button type='submit' className='py-2 px-4 bg-blue-500 text-white font-semibold rounded-lg shadow-md hover:bg-blue-700 focus:outline-hidden focus:ring-2 focus:ring-blue-400 focus:ring-opacity-75'>Submit</button>
        </form>
  )
}

export default LoginPage
