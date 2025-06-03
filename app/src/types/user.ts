export interface User {
    email: string
    first_name: string | null
    last_name: string | null
    is_superuser: boolean
    is_template_manager: boolean
}