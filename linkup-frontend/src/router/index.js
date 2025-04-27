import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '../pages/HomePage.vue'
import ProfilePage from '../pages/ProfilePage.vue'
import MyProfilePage from '../pages/MyProfilePage.vue'
import UpdateProfilePage from '../pages/UpdateProfilePage.vue'
import CreateThreadPage from '../pages/CreateThreadPage.vue'
import ThreadDetailPage from '../pages/ThreadDetailPage.vue'
import SearchPage from '../pages/SearchPage.vue'
import CommentPage from '../pages/CommentPage.vue'
import LoginPage from '../pages/LoginPage.vue'
import RegisterPage from '../pages/RegisterPage.vue'
import EmailVerifiedPage from '../pages/EmailVerifiedPage.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: HomePage,
        meta: { requiresAuth: true }
    },
    {
        path: '/me',
        name: 'my-profile',
        component: MyProfilePage,
        meta: { requiresAuth: true }
    },
    {
        path: '/me/edit',
        name: 'edit-profile',
        component: UpdateProfilePage,
        meta: { requiresAuth: true }
    },
    {
        path: '/profile/:id',
        name: 'profile',
        component: ProfilePage,
        props: true,
        meta: { requiresAuth: true }
    },
    { path: '/search', name: 'search', component: SearchPage, meta: { requiresAuth: true } },
    {
        path: '/thread/create',
        name: 'create-thread',
        component: CreateThreadPage,
        meta: { requiresAuth: true }
    },
    {
        path: '/thread/:id',
        name: 'thread-detail',
        component: ThreadDetailPage,
        props: true,
        meta: { requiresAuth: true }
    },
    {
        path: '/comments/:threadId',
        name: 'comments',
        component: CommentPage,
        props: true,
        meta: { requiresAuth: true }
    },
    {
        path: '/login',
        name: 'login',
        component: LoginPage,
        meta: { requiresGuest: true }
    },
    {
        path: '/register',
        name: 'register',
        component: RegisterPage,
        meta: { requiresGuest: true }
    },
    {
        path: '/email-verified',
        name: 'email-verified',
        component: EmailVerifiedPage,
        meta: { requiresGuest: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    if (to.meta.requiresAuth && !token) {
        return next({ name: 'login' })
    }
    if (to.meta.requiresGuest && token) {
        return next({ name: 'home' })
    }
    next()
})

export default router
