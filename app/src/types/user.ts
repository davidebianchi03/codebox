export interface User {
    email: string
    first_name: string
    last_name: string
    is_superuser: boolean
    is_template_manager: boolean
    last_login: string | null
}
