const m = window.m;
import i18n from "../services/i18n.js";
import HouseholdService from "../services/HouseholdService.js";

const PersonForm = {
    view: ({ attrs: { prefix, isHead, data, errors = {}, onRemove } }) => {
        return m("div.mb-8", [
            m("div.text-xl.font-semibold.mb-4.flex.justify-between.items-center", [
                m("span", isHead ? i18n.t('signup.hoh') : i18n.t('signup.othermembers')),
                !isHead && m("button.text-red-500.hover:text-red-700", {
                    onclick: onRemove
                }, "×")
            ]),
            
            // Name fields
            m("div.grid.grid-cols-1.md:grid-cols-2.gap-4", [
                m("div.form-group", [
                    m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.firstname')),
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "text",
                        name: `${prefix}FirstName`,
                        value: data[`${prefix}FirstName`] || '',
                        onchange: (e) => data[`${prefix}FirstName`] = e.target.value,
                        class: errors[`${prefix}FirstName`] ? 'border-red-500' : ''
                    }),
                    errors[`${prefix}FirstName`] && m("div.text-red-500.text-sm.mt-1", errors[`${prefix}FirstName`])
                ]),
                m("div.form-group", [
                    m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.lastname')),
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "text",
                        name: `${prefix}LastName`,
                        value: data[`${prefix}LastName`] || '',
                        onchange: (e) => data[`${prefix}LastName`] = e.target.value,
                        class: errors[`${prefix}LastName`] ? 'border-red-500' : ''
                    }),
                    errors[`${prefix}LastName`] && m("div.text-red-500.text-sm.mt-1", errors[`${prefix}LastName`])
                ])
            ]),
            
            // Date of Birth
            m("div.form-group.mt-4", [
                m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.dob')),
                m("div.grid.grid-cols-3.gap-4", [
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "text",
                        placeholder: i18n.t('misc.month'),
                        name: `${prefix}DobMonth`,
                        value: data[`${prefix}DobMonth`] || '',
                        onchange: (e) => data[`${prefix}DobMonth`] = e.target.value,
                        class: errors[`${prefix}DobMonth`] ? 'border-red-500' : ''
                    }),
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "text",
                        placeholder: i18n.t('misc.day'),
                        name: `${prefix}DobDay`,
                        value: data[`${prefix}DobDay`] || '',
                        onchange: (e) => data[`${prefix}DobDay`] = e.target.value,
                        class: errors[`${prefix}DobDay`] ? 'border-red-500' : ''
                    }),
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "text",
                        placeholder: i18n.t('misc.year'),
                        name: `${prefix}DobYear`,
                        value: data[`${prefix}DobYear`] || '',
                        onchange: (e) => data[`${prefix}DobYear`] = e.target.value,
                        class: errors[`${prefix}DobYear`] ? 'border-red-500' : ''
                    })
                ]),
                (errors[`${prefix}DobMonth`] || errors[`${prefix}DobDay`] || errors[`${prefix}DobYear`]) &&
                    m("div.text-red-500.text-sm.mt-1", i18n.t('misc.fieldrequired'))
            ]),
            
            // Contact Information (only for head of household)
            isHead && m("div.grid.grid-cols-1.md:grid-cols-2.gap-4.mt-4", [
                m("div.form-group", [
                    m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.email')),
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "email",
                        name: `${prefix}Email`,
                        value: data[`${prefix}Email`] || '',
                        onchange: (e) => data[`${prefix}Email`] = e.target.value
                    })
                ]),
                m("div.form-group", [
                    m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.phone')),
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "tel",
                        name: `${prefix}Phone`,
                        value: data[`${prefix}Phone`] || '',
                        onchange: (e) => data[`${prefix}Phone`] = e.target.value
                    })
                ])
            ]),
            
            // Address (only for head of household)
            isHead && m("div.mt-4", [
                m("div.form-group", [
                    m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.address')),
                    m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                        type: "text",
                        name: `${prefix}Street`,
                        value: data[`${prefix}Street`] || '',
                        onchange: (e) => data[`${prefix}Street`] = e.target.value
                    })
                ]),
                m("div.grid.grid-cols-2.gap-4.mt-4", [
                    m("div.form-group", [
                        m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.city')),
                        m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                            type: "text",
                            name: `${prefix}City`,
                            value: data[`${prefix}City`] || '',
                            onchange: (e) => data[`${prefix}City`] = e.target.value
                        })
                    ]),
                    m("div.form-group", [
                        m("label.block.text-sm.font-medium.text-gray-700", i18n.t('misc.zipcode')),
                        m("input.mt-1.block.w-full.rounded-md.border-gray-300.shadow-sm.focus:ring-blue-500.focus:border-blue-500", {
                            type: "text",
                            name: `${prefix}Zip`,
                            value: data[`${prefix}Zip`] || '',
                            onchange: (e) => data[`${prefix}Zip`] = e.target.value
                        })
                    ])
                ])
            ])
        ]);
    }
};

export const Signup = {
    data: {},
    errors: {},
    
    oninit: (vnode) => {
        vnode.state.showSuccess = false;
        vnode.state.members = []; // Track additional household members
    },
    
    addMember: (vnode) => {
        if (vnode.state.members.length < 5) {
            vnode.state.members.push({ id: Date.now() }); // Use timestamp as unique ID
        }
    },
    
    removeMember: (vnode, index) => {
        vnode.state.members.splice(index, 1);
    },
    
    submitForm: async (vnode) => {
        const household = {
            head: {
                firstName: vnode.state.data.hohFirstName,
                lastName: vnode.state.data.hohLastName,
                dob: `${vnode.state.data.hohDobYear}-${vnode.state.data.hohDobMonth}-${vnode.state.data.hohDobDay}`,
                email: vnode.state.data.hohEmail,
                phone: vnode.state.data.hohPhone,
                street: vnode.state.data.hohStreet,
                city: vnode.state.data.hohCity,
                postalCode: vnode.state.data.hohZip
            },
            members: vnode.state.members.map((member, index) => ({
                firstName: vnode.state.data[`person${index}FirstName`],
                lastName: vnode.state.data[`person${index}LastName`],
                dob: `${vnode.state.data[`person${index}DobYear`]}-${vnode.state.data[`person${index}DobMonth`]}-${vnode.state.data[`person${index}DobDay`]}`
            })).filter(member => member.firstName && member.lastName) // Only include members with at least first and last name
        };
        
        try {
            await HouseholdService.createHousehold(household);
            vnode.state.showSuccess = true;
            vnode.state.errors = {};
        } catch (error) {
            console.error('Error submitting form:', error);
            vnode.state.errors = error.response?.errors || { general: i18n.t('misc.error') };
        }
    },
    
    view: (vnode) => {
        if (vnode.state.showSuccess) {
            return m("div.container.mx-auto.px-4.py-8", [
                m("h1.text-3xl.font-bold.text-center.mb-4", i18n.t('misc.thankyou')),
                m("p.text-center", i18n.t('signup.success'))
            ]);
        }
        
        return m("div.container.mx-auto.px-4.py-8", [
            // Language Toggle
            m("div.flex.justify-end.mb-4", [
                m("button.text-blue-500.hover:text-blue-700.mr-2", {
                    onclick: () => i18n.setLanguage('en')
                }, "English"),
                m("span", "|"),
                m("button.text-blue-500.hover:text-blue-700.ml-2", {
                    onclick: () => i18n.setLanguage('es')
                }, "Español")
            ]),
            
            m("h1.text-3xl.font-bold.text-center.mb-4", i18n.t('signup.title')),
            m("p.text-center.mb-8", i18n.t('signup.intro')),
            
            m("form.max-w-4xl.mx-auto", {
                onsubmit: (e) => {
                    e.preventDefault();
                    vnode.state.submitForm(vnode);
                }
            }, [
                // Head of Household
                m(PersonForm, {
                    prefix: 'hoh',
                    isHead: true,
                    data: vnode.state.data,
                    errors: vnode.state.errors
                }),
                
                // Other Household Members Section
                m("div", [
                    vnode.state.members.length > 0 && m("h2.text-2xl.font-bold.mt-8.mb-4", i18n.t('signup.othermembers')),
                    
                    // Other Household Members
                    ...vnode.state.members.map((member, index) => 
                        m(PersonForm, {
                            prefix: `person${index}`,
                            isHead: false,
                            data: vnode.state.data,
                            errors: vnode.state.errors,
                            onRemove: () => vnode.state.removeMember(vnode, index)
                        })
                    )
                ]),
                
                // Add Member Button
                m("div.mt-4.text-center", 
                    vnode.state.members.length < 5 && m("button.text-blue-500.hover:text-blue-700.focus:outline-none[type=button]", {
                        onclick: () => vnode.state.addMember(vnode)
                    }, "Add Household Member")
                ),
                
                // Submit Button
                m("div.mt-8.text-center", [
                    vnode.state.errors.general && m("div.text-red-500.mb-4", vnode.state.errors.general),
                    m("button.bg-blue-500.text-white.px-6.py-2.rounded-md.hover:bg-blue-600.focus:outline-none.focus:ring-2.focus:ring-blue-500.focus:ring-offset-2[type=submit]",
                        i18n.t('misc.submit')
                    )
                ])
            ])
        ]);
    }
};

export default Signup;
