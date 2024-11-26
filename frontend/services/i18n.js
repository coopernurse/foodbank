const m = window.m;

const translations = {
    en: {
        'signup.title': 'Community Cupboard Sign-Up Form',
        'signup.intro': 'This information is helpful in providing our services. None of your information will be shared.',
        'signup.success': 'We have saved your information. Please ask for a shopping sheet from a staff member.',
        'signup.hoh': 'Head of Household',
        'signup.othermembers': 'Others Living in the Household',
        'misc.firstname': 'First Name',
        'misc.lastname': 'Last Name',
        'misc.address': 'Address',
        'misc.city': 'City',
        'misc.zipcode': 'ZIP Code',
        'misc.email': 'Email',
        'misc.phone': 'Phone',
        'misc.gender': 'Gender',
        'misc.male': 'Male',
        'misc.female': 'Female',
        'misc.prefernottosay': 'Prefer not to say',
        'misc.dob': 'Date of Birth',
        'misc.month': 'Month',
        'misc.day': 'Day',
        'misc.year': 'Year',
        'misc.primarylang': 'Primary Language',
        'misc.english': 'English',
        'misc.spanish': 'Spanish',
        'misc.other': 'Other',
        'misc.relationship': 'Relationship',
        'misc.child': 'Child',
        'misc.fieldrequired': 'This field is required',
        'misc.submit': 'Submit',
        'misc.thankyou': 'Thank You',
        'misc.error': 'An error occurred',
        'misc.select': 'Select...',
        'misc.race': 'Race',
        'misc.race.white': 'White/Anglo',
        'misc.race.latino': 'Latina/Latino',
        'misc.race.black': 'Black/African American',
        'misc.race.asian': 'Asian',
    },
    es: {
        'signup.title': 'Formulario de Inscripción de Community Cupboard',
        'signup.intro': 'Esta información es útil para proporcionar nuestros servicios. Su información no será compartida.',
        'signup.success': 'Hemos guardado su información. Por favor, solicite una hoja de compras a un miembro del personal.',
        'signup.hoh': 'Cabeza de Familia',
        'signup.othermembers': 'Otros Miembros del Hogar',
        'misc.firstname': 'Nombre',
        'misc.lastname': 'Apellido',
        'misc.address': 'Dirección',
        'misc.city': 'Ciudad',
        'misc.zipcode': 'Código Postal',
        'misc.email': 'Correo Electrónico',
        'misc.phone': 'Teléfono',
        'misc.gender': 'Género',
        'misc.male': 'Masculino',
        'misc.female': 'Femenino',
        'misc.prefernottosay': 'Prefiero no decir',
        'misc.dob': 'Fecha de Nacimiento',
        'misc.month': 'Mes',
        'misc.day': 'Día',
        'misc.year': 'Año',
        'misc.primarylang': 'Idioma Principal',
        'misc.english': 'Inglés',
        'misc.spanish': 'Español',
        'misc.other': 'Otro',
        'misc.relationship': 'Relación',
        'misc.child': 'Hijo/a',
        'misc.fieldrequired': 'Este campo es obligatorio',
        'misc.submit': 'Enviar',
        'misc.thankyou': 'Gracias',
        'misc.error': 'Ocurrió un error',
        'misc.select': 'Seleccionar...',
        'misc.race': 'Raza',
        'misc.race.white': 'Blanco/Anglo',
        'misc.race.latino': 'Latina/Latino',
        'misc.race.black': 'Negro/Afroamericano',
        'misc.race.asian': 'Asiático',
    }
};

const i18n = {
    currentLang: 'en',
    
    setLanguage(lang) {
        this.currentLang = lang;
        m.redraw();
    },
    
    t(key) {
        return translations[this.currentLang][key] || key;
    }
};

export default i18n;
