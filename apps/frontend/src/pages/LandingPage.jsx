import React from 'react';
import { Link } from 'react-router-dom';

/* ── Icon helpers ─────────────────────────────────── */
const IconStar = () => (
  <svg width="12" height="12" viewBox="0 0 512 512" fill="white">
    <path d="M361.5 1.2c5 2.1 8.6 6.6 9.6 11.9L391 121l107.9 19.8c5.3 1 9.8 4.6 11.9 9.6s1.5 10.7-1.6 15.2L446.9 256l62.3 90.3c3.1 4.5 3.7 10.2 1.6 15.2s-6.6 8.6-11.9 9.6L391 391 371.1 498.9c-1 5.3-4.6 9.8-9.6 11.9s-10.7 1.5-15.2-1.6L256 446.9l-90.3 62.3c-4.5 3.1-10.2 3.7-15.2 1.6s-8.6-6.6-9.6-11.9L121 391 13.1 371.1c-5.3-1-9.8-4.6-11.9-9.6s-1.5-10.7 1.6-15.2L65.1 256 2.8 165.7c-3.1-4.5-3.7-10.2-1.6-15.2s6.6-8.6 11.9-9.6L121 121 140.9 13.1c1-5.3 4.6-9.8 9.6-11.9s10.7-1.5 15.2 1.6L256 65.1 346.3 2.8c4.5-3.1 10.2-3.7 15.2-1.6z"/>
  </svg>
);

const IconCheck = ({ color = '#00a650', size = 14 }) => (
  <svg width={size} height={size} viewBox="0 0 512 512" fill={color}>
    <path d="M470.6 105.4c12.5 12.5 12.5 32.8 0 45.3l-256 256c-12.5 12.5-32.8 12.5-45.3 0l-128-128c-12.5-12.5-12.5-32.8 0-45.3s32.8-12.5 45.3 0L192 338.7 425.4 105.4c12.5-12.5 32.8-12.5 45.3 0z"/>
  </svg>
);

/* ── Dot pattern SVG ──────────────────────────────── */
const patternBg = `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.04'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`;

/* ══════════════════════════════════════════════════ */
const LandingPage = () => (
  <div className="font-sans text-gray-100 bg-gray-900 overflow-x-hidden">

    {/* ── Nav ──────────────────────────────────────── */}
    <nav style={{ position:'fixed', top:0, left:0, right:0, zIndex:100, background:'rgba(17,24,39,.93)', backdropFilter:'blur(12px)', borderBottom:'1px solid #374151', height:64, display:'flex', alignItems:'center' }}>
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px', width:'100%', display:'flex', alignItems:'center', justifyContent:'space-between' }}>
        <a href="#" style={{ display:'flex', alignItems:'center', gap:12, textDecoration:'none' }}>
          <img src="/logo64.png" alt="Niloft" style={{ width:36, height:36, borderRadius:'50%' }} />
          <span style={{ fontWeight:700, fontSize:'1.1rem', color:'#f3f4f6' }}>Niloft</span>
        </a>
        <div style={{ display:'flex', alignItems:'center', gap:24 }}>
          <a href="#features" style={{ fontSize:'0.875rem', fontWeight:500, color:'#9ca3af', textDecoration:'none' }}>Funciones</a>
          <a href="#how" style={{ fontSize:'0.875rem', fontWeight:500, color:'#9ca3af', textDecoration:'none' }}>Cómo funciona</a>
          <a href="#pricing" style={{ fontSize:'0.875rem', fontWeight:500, color:'#9ca3af', textDecoration:'none' }}>Planes</a>
          <Link to="/login" style={{ fontSize:'0.875rem', fontWeight:500, color:'#9ca3af', textDecoration:'none' }}>Iniciar sesión</Link>
          <Link to="/register" style={{ background:'#009ee3', color:'#fff', padding:'8px 20px', borderRadius:8, fontWeight:600, fontSize:'0.875rem', textDecoration:'none' }}>Empezá gratis</Link>
        </div>
      </div>
    </nav>

    {/* ── Hero ─────────────────────────────────────── */}
    <section style={{ minHeight:'100vh', display:'flex', alignItems:'center', background:'linear-gradient(135deg, #009ee3 0%, #00a650 100%)', paddingTop:64, position:'relative', overflow:'hidden' }}>
      <div style={{ position:'absolute', inset:0, backgroundImage: patternBg }} />
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px', width:'100%' }}>
        <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:60, alignItems:'center', padding:'80px 0', position:'relative', zIndex:1 }}>
          {/* Text */}
          <div>
            <div style={{ display:'inline-flex', alignItems:'center', gap:6, background:'rgba(255,255,255,.15)', color:'#fff', padding:'6px 14px', borderRadius:20, fontSize:'0.78rem', fontWeight:500, border:'1px solid rgba(255,255,255,.25)', marginBottom:20 }}>
              <IconStar /> Tu asistente financiero personal
            </div>
            <h1 style={{ fontSize:'3.2rem', fontWeight:800, color:'#fff', lineHeight:1.15, marginBottom:20 }}>
              Tomá el control<br />
              <span style={{ opacity:.85 }}>de tus finanzas</span>
            </h1>
            <p style={{ fontSize:'1.1rem', color:'rgba(255,255,255,.88)', lineHeight:1.65, marginBottom:36, maxWidth:480 }}>
              Registrá gastos e ingresos, creá presupuestos y recibí análisis con IA — todo en un solo lugar.
            </p>
            <div style={{ display:'flex', gap:12, flexWrap:'wrap' }}>
              <Link to="/register" style={{ background:'#fff', color:'#009ee3', padding:'14px 28px', borderRadius:10, fontWeight:700, fontSize:'0.95rem', textDecoration:'none', boxShadow:'0 4px 16px rgba(0,0,0,.12)' }}>
                Empezá gratis
              </Link>
              <a href="#how" style={{ background:'rgba(255,255,255,.15)', color:'#fff', padding:'14px 28px', borderRadius:10, fontWeight:600, fontSize:'0.95rem', textDecoration:'none', border:'1px solid rgba(255,255,255,.3)' }}>
                Ver cómo funciona
              </a>
            </div>
            <p style={{ marginTop:16, fontSize:'0.78rem', color:'rgba(255,255,255,.65)', display:'flex', alignItems:'center', gap:6 }}>
              <IconCheck color="rgba(255,255,255,.65)" size={12} />
              Demo sin tarjeta de crédito · Cancelá cuando quieras
            </p>
          </div>

          {/* App mockup */}
          <div style={{ position:'relative' }}>
            <div style={{ background:'#fff', borderRadius:16, boxShadow:'0 24px 80px rgba(0,0,0,.25)', overflow:'hidden', border:'1px solid rgba(255,255,255,.3)' }}>
              {/* Window bar */}
              <div style={{ background:'#f3f4f6', padding:'10px 14px', display:'flex', alignItems:'center', gap:6, borderBottom:'1px solid #e5e7eb' }}>
                <div style={{ width:10, height:10, borderRadius:'50%', background:'#ff5f57' }} />
                <div style={{ width:10, height:10, borderRadius:'50%', background:'#ffbd2e' }} />
                <div style={{ width:10, height:10, borderRadius:'50%', background:'#28c840' }} />
                <div style={{ flex:1, background:'#e5e7eb', borderRadius:4, height:16, marginLeft:8, display:'flex', alignItems:'center', padding:'0 8px' }}>
                  <span style={{ fontSize:10, color:'#9ca3af' }}>financial.niloft.com</span>
                </div>
              </div>
              {/* Mini app */}
              <div style={{ display:'flex', height:340, overflow:'hidden', background:'#f9fafb', fontFamily:'Inter, sans-serif', fontSize:11 }}>
                {/* Sidebar */}
                <div style={{ width:52, background:'#fff', borderRight:'1px solid #e5e7eb', display:'flex', flexDirection:'column', alignItems:'center', padding:'12px 8px', gap:8 }}>
                  <img src="/logo64.png" width={28} height={28} style={{ borderRadius:'50%', marginBottom:4 }} alt="" />
                  {[['#eff6ff','#2563eb'],['#f9fafb','#9ca3af'],['#f9fafb','#9ca3af'],['#f9fafb','#9ca3af']].map(([bg, fill], i) => (
                    <div key={i} style={{ width:32, height:32, borderRadius:8, background:bg, display:'flex', alignItems:'center', justifyContent:'center' }}>
                      <div style={{ width:14, height:14, borderRadius:3, background:fill, opacity:.6 }} />
                    </div>
                  ))}
                </div>
                {/* Main */}
                <div style={{ flex:1 }}>
                  <div style={{ background:'#fff', borderBottom:'1px solid #e5e7eb', padding:'8px 14px', display:'flex', justifyContent:'space-between', alignItems:'center' }}>
                    <span style={{ fontSize:11, fontWeight:600, color:'#111827' }}>Cuentas · Abril 2026</span>
                    <div style={{ width:22, height:22, borderRadius:'50%', background:'linear-gradient(135deg,#009ee3,#00a650)', display:'flex', alignItems:'center', justifyContent:'center', color:'#fff', fontSize:9, fontWeight:700 }}>M</div>
                  </div>
                  <div style={{ padding:14, display:'flex', flexDirection:'column', gap:10 }}>
                    <div style={{ display:'grid', gridTemplateColumns:'repeat(4,1fr)', gap:8 }}>
                      {[['Balance','$71.300','▲ Positivo','#00a650'],['Ingresos','$111.300','3 registros','#111827'],['Gastos','$40.000','12 gastos','#111827'],['Pendientes','4','Por pagar','#ff6900']].map(([label, val, sub, color]) => (
                        <div key={label} style={{ background:'#fff', borderRadius:8, border:'1px solid #e5e7eb', padding:'10px 12px' }}>
                          <div style={{ fontSize:9, color:'#6b7280', marginBottom:3 }}>{label}</div>
                          <div style={{ fontSize:13, fontWeight:700, color }}>{val}</div>
                          <div style={{ fontSize:8, color, marginTop:2 }}>{sub}</div>
                        </div>
                      ))}
                    </div>
                    <div style={{ background:'#fff', borderRadius:8, border:'1px solid #e5e7eb', padding:'10px 12px' }}>
                      <div style={{ display:'flex', justifyContent:'space-between', marginBottom:8 }}>
                        <span style={{ fontSize:10, fontWeight:600, color:'#111827' }}>Transacciones recientes</span>
                        <span style={{ fontSize:9, color:'#009ee3' }}>Ver todas →</span>
                      </div>
                      {[['#3b82f6','Supermercado','-$4.850','#111827'],['#00a650','Sueldo neto','+$85.000','#00a650'],['#e53e3e','Tarjeta Visa','-$8.400','#111827']].map(([dot, name, amount, amtColor], i, arr) => (
                        <div key={name} style={{ display:'flex', justifyContent:'space-between', alignItems:'center', padding:'5px 0', borderBottom: i < arr.length-1 ? '1px solid #f3f4f6' : 'none' }}>
                          <div style={{ display:'flex', alignItems:'center', gap:6 }}>
                            <div style={{ width:6, height:6, borderRadius:'50%', background:dot }} />
                            <span style={{ fontSize:10, color:'#374151' }}>{name}</span>
                          </div>
                          <span style={{ fontSize:10, fontWeight:600, color:amtColor }}>{amount}</span>
                        </div>
                      ))}
                    </div>
                    <div style={{ background:'#fff', borderRadius:8, border:'1px solid #e5e7eb', padding:'10px 12px' }}>
                      <div style={{ fontSize:10, fontWeight:600, color:'#111827', marginBottom:8 }}>Ingresos vs. Gastos</div>
                      <div style={{ display:'flex', alignItems:'flex-end', gap:6, height:36 }}>
                        {[['#dcfce7',100],['#fee2e2',36],['#dcfce7',88],['#fee2e2',30],['#dbeafe',72],['#fee2e2',38]].map(([bg, h], i) => (
                          <div key={i} style={{ flex:1, background:bg, borderRadius:'3px 3px 0 0', height:`${h}%` }} />
                        ))}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            {/* Floating badge */}
            <div style={{ position:'absolute', bottom:-16, right:-16, background:'#1f2937', borderRadius:12, padding:'10px 14px', boxShadow:'0 8px 24px rgba(0,0,0,.4)', display:'flex', alignItems:'center', gap:8, border:'1px solid #374151' }}>
              <div style={{ width:32, height:32, borderRadius:'50%', background:'rgba(0,166,80,.2)', display:'flex', alignItems:'center', justifyContent:'center' }}>
                <IconCheck />
              </div>
              <div>
                <div style={{ fontSize:11, fontWeight:700, color:'#f3f4f6' }}>Presupuesto cumplido</div>
                <div style={{ fontSize:10, color:'#9ca3af' }}>Ahorraste $8.200 este mes</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    {/* ── Stats ────────────────────────────────────── */}
    <div style={{ background:'#1f2937', borderTop:'1px solid #374151', borderBottom:'1px solid #374151', padding:'40px 0' }}>
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px' }}>
        <div style={{ display:'grid', gridTemplateColumns:'repeat(3,1fr)' }}>
          {[
            ['IA incorporada', 'Único desde Argentina', '#009ee3'],
            ['Tips exclusivos', 'La inteligencia artificial a tu servicio', '#f3f4f6'],
            ['Intuitivo', 'Finanzas sofisticadas sin dificultad', '#f3f4f6'],
          ].map(([num, label, color], i, arr) => (
            <div key={num} style={{ textAlign:'center', padding:24, borderRight: i < arr.length-1 ? '1px solid #374151' : 'none' }}>
              <div style={{ fontSize:'2rem', fontWeight:800, color }}>{num}</div>
              <div style={{ fontSize:'0.875rem', color:'#9ca3af', marginTop:4 }}>{label}</div>
            </div>
          ))}
        </div>
      </div>
    </div>

    {/* ── Features ─────────────────────────────────── */}
    <section id="features" style={{ background:'#111827', padding:'80px 0' }}>
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px' }}>
        <div style={{ fontSize:'0.78rem', fontWeight:600, color:'#009ee3', textTransform:'uppercase', letterSpacing:'.08em', marginBottom:12 }}>Funciones</div>
        <h2 style={{ fontSize:'2.2rem', fontWeight:800, color:'#f3f4f6', marginBottom:14, lineHeight:1.2 }}>
          Todo lo que necesitás<br />para ordenar tu dinero
        </h2>
        <p style={{ fontSize:'1rem', color:'#9ca3af', maxWidth:560, lineHeight:1.65 }}>
          Financial Niloft no es solo un registro de gastos. Es un sistema completo con educación financiera personalizada para entender, planificar y mejorar tus finanzas.
        </p>
        <div style={{ display:'grid', gridTemplateColumns:'repeat(3,1fr)', gap:24, marginTop:56 }}>
          {[
            { bg:'rgba(37,99,235,.15)', fill:'#60a5fa', title:'Gastos e ingresos', desc:'Registrá todos tus movimientos en segundos. Categorizá, marcá como pagado y buscá cualquier transacción al instante.' },
            { bg:'rgba(0,166,80,.15)', fill:'#00a650', title:'Presupuestos inteligentes', desc:'Creá límites de gasto por categoría. Niloft te avisa cuando te acercás al tope para que nunca te sorprendas a fin de mes.' },
            { bg:'rgba(124,58,237,.15)', fill:'#a78bfa', title:'Asesor con IA', desc:'Analizá tu salud financiera con inteligencia artificial. Recibí consejos personalizados, simulaciones y un score financiero propio.' },
            { bg:'rgba(194,65,12,.15)', fill:'#fb923c', title:'Objetivos de ahorro', desc:'Definí metas con montos y fechas. Realizá depósitos parciales y seguí tu progreso visualmente hasta alcanzar el objetivo.' },
            { bg:'rgba(0,166,80,.15)', fill:'#00a650', title:'Transacciones recurrentes', desc:'Programá pagos fijos como suscripciones, alquiler o cuotas. Niloft los ejecuta automáticamente y te avisa cuando vencen.' },
            { bg:'rgba(37,99,235,.15)', fill:'#60a5fa', title:'Reportes detallados', desc:'Visualizá tendencias mensuales, distribución por categorías y evolución de tu patrimonio. Exportá tus datos cuando quieras.' },
          ].map(({ bg, fill, title, desc }) => (
            <div key={title} style={{ padding:28, borderRadius:14, border:'1px solid #374151', background:'#1f2937', transition:'all .2s' }}
              onMouseEnter={e => { e.currentTarget.style.boxShadow='0 8px 32px rgba(0,0,0,.3)'; e.currentTarget.style.transform='translateY(-2px)'; e.currentTarget.style.borderColor='#3b82f6'; }}
              onMouseLeave={e => { e.currentTarget.style.boxShadow='none'; e.currentTarget.style.transform=''; e.currentTarget.style.borderColor='#374151'; }}>
              <div style={{ width:44, height:44, borderRadius:12, background:bg, display:'flex', alignItems:'center', justifyContent:'center', marginBottom:16 }}>
                <div style={{ width:20, height:20, borderRadius:4, background:fill, opacity:.9 }} />
              </div>
              <div style={{ fontSize:'1rem', fontWeight:700, color:'#f3f4f6', marginBottom:8 }}>{title}</div>
              <div style={{ fontSize:'0.875rem', color:'#9ca3af', lineHeight:1.6 }}>{desc}</div>
            </div>
          ))}
        </div>
      </div>
    </section>

    {/* ── How it works ─────────────────────────────── */}
    <section id="how" style={{ background:'#1f2937', padding:'80px 0' }}>
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px' }}>
        <div style={{ textAlign:'center', marginBottom:0 }}>
          <div style={{ fontSize:'0.78rem', fontWeight:600, color:'#009ee3', textTransform:'uppercase', letterSpacing:'.08em', marginBottom:12 }}>Cómo funciona</div>
          <h2 style={{ fontSize:'2.2rem', fontWeight:800, color:'#f3f4f6', marginBottom:14, lineHeight:1.2 }}>Tres pasos y ya estás</h2>
          <p style={{ fontSize:'1rem', color:'#9ca3af', maxWidth:480, margin:'0 auto', lineHeight:1.65 }}>
            Configurá tu cuenta en minutos y empezá a ver el impacto de inmediato.
          </p>
        </div>
        <div style={{ display:'grid', gridTemplateColumns:'repeat(3,1fr)', gap:32, marginTop:56, position:'relative' }}>
          <div style={{ position:'absolute', top:28, left:'calc(16.66% + 28px)', right:'calc(16.66% + 28px)', height:2, background:'linear-gradient(90deg,#009ee3,#00a650)', zIndex:0 }} />
          {[
            ['1', 'Creá tu cuenta gratis', 'Registrate con tu email en menos de un minuto. Sin configuraciones complejas ni datos bancarios requeridos.'],
            ['2', 'Cargá tus movimientos', 'Ingresá tus gastos e ingresos del mes. Podés crear categorías propias o usar las predefinidas de Niloft.'],
            ['3', 'Tomá decisiones informadas', 'El dashboard y el asesor IA te muestran exactamente dónde está tu dinero y cómo mejorar tu salud financiera.'],
          ].map(([num, title, desc]) => (
            <div key={num} style={{ textAlign:'center', position:'relative', zIndex:1 }}>
              <div style={{ width:56, height:56, borderRadius:'50%', display:'flex', alignItems:'center', justifyContent:'center', fontWeight:800, fontSize:'1.1rem', color:'#fff', margin:'0 auto 18px', background:'linear-gradient(135deg,#009ee3,#00a650)', boxShadow:'0 4px 12px rgba(0,158,227,.3)' }}>{num}</div>
              <div style={{ fontSize:'1.05rem', fontWeight:700, color:'#f3f4f6', marginBottom:8 }}>{title}</div>
              <div style={{ fontSize:'0.875rem', color:'#9ca3af', lineHeight:1.6 }}>{desc}</div>
            </div>
          ))}
        </div>
      </div>
    </section>

    {/* ── Pricing ──────────────────────────────────── */}
    <section id="pricing" style={{ background:'#111827', padding:'80px 0' }}>
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px' }}>
        <div style={{ textAlign:'center' }}>
          <div style={{ fontSize:'0.78rem', fontWeight:600, color:'#009ee3', textTransform:'uppercase', letterSpacing:'.08em', marginBottom:12 }}>Planes</div>
          <h2 style={{ fontSize:'2.2rem', fontWeight:800, color:'#f3f4f6', marginBottom:14, lineHeight:1.2 }}>Simple, transparente,<br />sin sorpresas</h2>
          <p style={{ fontSize:'1rem', color:'#9ca3af', maxWidth:480, margin:'0 auto', lineHeight:1.65 }}>
            Empezá gratis y desbloqueá más funciones a medida que avanzás.
          </p>
        </div>
        <div style={{ display:'grid', gridTemplateColumns:'repeat(3,1fr)', gap:24, marginTop:56 }}>
          {/* Free */}
          <PricingCard
            name="Básico"
            tagline="Para empezar a ordenarte"
            price="Gratis"
            period="Para siempre"
            features={['Registro de gastos e ingresos', 'Categorías personalizadas', 'Dashboard mensual', 'Historial de 3 meses']}
            disabled={['Presupuestos inteligentes', 'Objetivos de ahorro', 'Asesor con IA']}
            btnLabel="Empezá gratis"
            btnStyle="ghost"
          />
          {/* Pro */}
          <PricingCard
            featured
            badge="Más popular"
            name="Pro"
            tagline="Para gestionar en serio"
            price="$9.990"
            period="Facturación mensual · Cancelá cuando quieras"
            features={['Todo lo del plan Básico', 'Presupuestos por categoría', 'Objetivos de ahorro', 'Transacciones recurrentes', 'Reportes avanzados', 'Historial ilimitado']}
            disabled={['Asesor con IA']}
            btnLabel="Empezá 14 días gratis"
            btnStyle="primary"
          />
          {/* Premium */}
          <PricingCard
            name="Premium"
            tagline="Con inteligencia artificial"
            price="$15.490"
            period="Facturación mensual · Cancelá cuando quieras"
            features={['Todo lo del plan Pro', 'Asesor financiero con IA', 'Score financiero personal', 'Análisis predictivo', 'Simulador "¿Lo puedo comprar?"', 'Soporte prioritario', 'Multi-miembros / familia']}
            disabled={[]}
            btnLabel="Empezá 14 días gratis"
            btnStyle="outline"
          />
        </div>
      </div>
    </section>

    {/* ── CTA Banner ───────────────────────────────── */}
    <section style={{ background:'linear-gradient(135deg,#009ee3 0%,#00a650 100%)', position:'relative', overflow:'hidden', padding:'80px 0' }}>
      <div style={{ position:'absolute', inset:0, backgroundImage: patternBg }} />
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px', textAlign:'center', position:'relative', zIndex:1 }}>
        <h2 style={{ fontSize:'2.4rem', fontWeight:800, color:'#fff', marginBottom:14 }}>¿Listo para tomar el control?</h2>
        <p style={{ fontSize:'1.05rem', color:'rgba(255,255,255,.85)', marginBottom:36, maxWidth:480, marginLeft:'auto', marginRight:'auto' }}>
          Unite a miles de personas que ya organizan sus finanzas con Niloft. Gratis, en español, sin complicaciones.
        </p>
        <Link to="/register" style={{ background:'#fff', color:'#009ee3', padding:'14px 32px', borderRadius:10, fontWeight:700, fontSize:'1rem', textDecoration:'none', display:'inline-block', boxShadow:'0 4px 16px rgba(0,0,0,.15)' }}>
          Empezá gratis hoy
        </Link>
        <p style={{ marginTop:14, fontSize:'0.8rem', color:'rgba(255,255,255,.65)' }}>Sin tarjeta de crédito · Configuración en 2 minutos</p>
      </div>
    </section>

    {/* ── Footer ───────────────────────────────────── */}
    <footer style={{ background:'#0d1117', color:'#9ca3af', padding:'56px 0 32px' }}>
      <div style={{ maxWidth:1120, margin:'0 auto', padding:'0 24px' }}>
        <div style={{ display:'grid', gridTemplateColumns:'2fr 1fr 1fr 1fr', gap:48, marginBottom:48 }}>
          <div>
            <div style={{ display:'flex', alignItems:'center', gap:10, marginBottom:12 }}>
              <img src="/logo64.png" alt="Niloft" style={{ width:36, height:36, borderRadius:'50%' }} />
              <span style={{ fontWeight:700, fontSize:'1.1rem', color:'#fff' }}>Niloft</span>
            </div>
            <p style={{ fontSize:'0.85rem', lineHeight:1.65 }}>Tu asistente financiero personal. Organizá tus finanzas, alcanzá tus metas y tomá mejores decisiones económicas.</p>
          </div>
          {[
            ['Producto', ['Funciones', 'Precios', 'Novedades', 'Roadmap']],
            ['Soporte', ['Centro de ayuda', 'Contacto', 'Estado del servicio']],
            ['Legal', ['Términos de uso', 'Privacidad', 'Cookies']],
          ].map(([title, links]) => (
            <div key={title}>
              <div style={{ fontSize:'0.82rem', fontWeight:600, color:'#fff', marginBottom:16, textTransform:'uppercase', letterSpacing:'.06em' }}>{title}</div>
              <div style={{ display:'flex', flexDirection:'column', gap:10 }}>
                {links.map(link => (
                  <a key={link} href="#" style={{ fontSize:'0.85rem', color:'#9ca3af', textDecoration:'none' }}>{link}</a>
                ))}
              </div>
            </div>
          ))}
        </div>
        <div style={{ borderTop:'1px solid #374151', paddingTop:24, display:'flex', justifyContent:'space-between', alignItems:'center' }}>
          <span style={{ fontSize:'0.8rem' }}>© 2026 Niloft. Hecho con ♥ en Argentina.</span>
          <div style={{ display:'flex', gap:20 }}>
            <a href="#" style={{ fontSize:'0.8rem', color:'#9ca3af', textDecoration:'none' }}>Términos</a>
            <a href="#" style={{ fontSize:'0.8rem', color:'#9ca3af', textDecoration:'none' }}>Privacidad</a>
          </div>
        </div>
      </div>
    </footer>
  </div>
);

/* ── Pricing card sub-component ──────────────────── */
const PricingCard = ({ featured, badge, name, tagline, price, period, features, disabled, btnLabel, btnStyle }) => (
  <div style={{
    padding:'32px 28px', borderRadius:16, position:'relative', background:'#1f2937',
    border: featured ? '2px solid #009ee3' : '1px solid #374151',
    boxShadow: featured ? '0 8px 32px rgba(0,158,227,.2)' : 'none',
  }}>
    {badge && (
      <div style={{ position:'absolute', top:-12, left:'50%', transform:'translateX(-50%)', background:'linear-gradient(135deg,#009ee3,#00a650)', color:'#fff', padding:'4px 14px', borderRadius:20, fontSize:'0.72rem', fontWeight:700, whiteSpace:'nowrap' }}>
        {badge}
      </div>
    )}
    <div style={{ fontSize:'1rem', fontWeight:700, color:'#f3f4f6', marginBottom:6 }}>{name}</div>
    <div style={{ fontSize:'0.82rem', color:'#9ca3af', marginBottom:12 }}>{tagline}</div>
    <div style={{ fontSize:'2.4rem', fontWeight:800, color:'#f3f4f6', margin:'12px 0 4px' }}>
      {price}<span style={{ fontSize:'1rem', fontWeight:500, color:'#9ca3af' }}>{price !== 'Gratis' ? '/mes' : ''}</span>
    </div>
    <div style={{ fontSize:'0.8rem', color:'#6b7280', marginBottom:20 }}>{period}</div>
    <ul style={{ listStyle:'none', marginBottom:28, display:'flex', flexDirection:'column', gap:10 }}>
      {features.map(f => (
        <li key={f} style={{ display:'flex', alignItems:'center', gap:8, fontSize:'0.85rem', color:'#d1d5db' }}>
          <span style={{ width:16, height:16, borderRadius:'50%', background:'rgba(0,166,80,.2)', display:'inline-flex', alignItems:'center', justifyContent:'center', flexShrink:0 }}>
            <IconCheck size={10} />
          </span>
          {f}
        </li>
      ))}
      {disabled.map(f => (
        <li key={f} style={{ display:'flex', alignItems:'center', gap:8, fontSize:'0.85rem', color:'#6b7280' }}>
          <span style={{ width:16, height:16, borderRadius:'50%', background:'#374151', display:'inline-flex', alignItems:'center', justifyContent:'center', flexShrink:0 }}>
            <svg width="10" height="10" viewBox="0 0 10 10"><path d="M2 2l6 6M8 2l-6 6" stroke="#6b7280" strokeWidth="1.5" strokeLinecap="round"/></svg>
          </span>
          {f}
        </li>
      ))}
    </ul>
    <Link
      to="/register"
      style={{
        display:'block', textAlign:'center', padding:12, borderRadius:10, fontWeight:600, fontSize:'0.9rem', textDecoration:'none', transition:'all .15s',
        ...(btnStyle === 'primary' ? { background:'#009ee3', color:'#fff' } :
           btnStyle === 'outline' ? { background:'transparent', color:'#009ee3', border:'2px solid #009ee3' } :
           { background:'#374151', color:'#d1d5db' }),
      }}
    >
      {btnLabel}
    </Link>
  </div>
);

export default LandingPage;
