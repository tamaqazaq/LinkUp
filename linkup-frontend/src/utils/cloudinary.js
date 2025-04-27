export async function uploadToCloudinary(file, userId) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('upload_preset', import.meta.env.VITE_CLOUDINARY_UPLOAD_PRESET)
    formData.append('public_id', `avatar_${userId}`) // 👈 перезапишет аватар
    formData.append('overwrite', true)

    const res = await fetch(
        `https://api.cloudinary.com/v1_1/${import.meta.env.VITE_CLOUDINARY_CLOUD_NAME}/upload`,
        { method: 'POST', body: formData }
    )

    const data = await res.json()
    return data.secure_url + `?t=${Date.now()}`
}
