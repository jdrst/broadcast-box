import { useContext, useEffect } from 'react';
import { AuthContext } from '../login'
import { Link, Outlet, useNavigate } from 'react-router-dom'
import { CinemaModeContext } from '../player';


const Header = () => {
  const { cinemaMode } = useContext(CinemaModeContext);
  const navbarEnabled = !cinemaMode;
  const navigate = useNavigate();
  const [{username, authenticated, pending}, login, logout] = useContext(AuthContext);
  
  useEffect(() => {
    if (!authenticated && !pending) {
      const url = new URL(document.location)
      const from = url.searchParams.get("redirectUrl") ?? `${url.pathname}${url.search}`
      navigate(`/login?redirectUrl=${encodeURIComponent(from)}`, { replace: true })
    }
  }, [authenticated, pending])
  
  const handleLogoutClick = () => {
    logout().then(resp => {
      if (resp.ok) {
        navigate('/')
      }
    })
  }
  
  return (
    <div>
      {navbarEnabled && (
        <nav className='bg-gray-800 p-2 mt-0 fixed w-full z-10 top-0'>
          <div className='container mx-auto flex flex-wrap items-center'>
            <div className='flex flex-1 text-white font-extrabold'>
              <Link to="/" className='font-light leading-tight text-2xl'>
                Broadcast Box
              </Link>
              {authenticated &&
                <>
                  Hello {username}
                  <button onClick={handleLogoutClick} className='ml-10 py-2 px-4 bg-blue-500 text-white font-semibold rounded-lg shadow-md hover:bg-blue-700 focus:outline-hidden focus:ring-2 focus:ring-blue-400 focus:ring-opacity-75'>Logout</button>
                </>
              }
            </div>
          </div>
        </nav>
      )}

      <main className={`${navbarEnabled && "pt-20 md:pt-24"}`}>
        <Outlet />
      </main>

      <footer className="mx-auto px-2 container py-6">
        <ul className="flex items-center justify-center mt-3 text-sm:mt-0 space-x-4">
          <li>
            <a href="https://github.com/Glimesh/broadcast-box" className="hover:underline">GitHub</a>
          </li>
          <li>
            <a href="https://pion.ly" className="hover:underline">Pion</a>
          </li>
          <li>
            <a href="https://glimesh.tv" className="hover:underline">Glimesh</a>
          </li>
        </ul>
      </footer>

    </div>
  )
}

export default Header
