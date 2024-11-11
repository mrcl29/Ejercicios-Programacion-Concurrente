with Text_Io;
use  Text_Io;
with Ada.Numerics.Discrete_Random;
with Ada.Numerics.Float_Random;
with productors_consumidors;
use productors_consumidors;

procedure maquinarefrescs is

   -- Declaracio per a nombres aleatoris
   queden: integer;
   subtype rang_tipus is Integer range 0 .. 10;
   package aleatori is new Ada.Numerics.Discrete_Random (rang_tipus); -- Necessitem un nombre aleatori entre [0,10]
   Gen: aleatori.Generator;  -- Generador de nombres aleatoris
   Monitor: Monitor_Productor_Consumidor;

   type cadena_variable is access all String;
   final: boolean := False;

   -- Definicio de les tasques
   task type tasca_client (nom_client: cadena_variable; n_reposadors: integer);
   task type tasca_reposador (n_reposador: Integer; clients: integer);  -- Tasca que representa un reposador

   -- Cos de la tasca client
   task body tasca_client is
      nom: constant String := nom_client.all;
      n_consumicions: integer;
   begin
      -- Si no hi ha reposadors, marxem i tenim 0 consumicions
      if(n_reposadors = 0) then
         Put_Line(nom& " diu: No hi ha reposadors a la maquina, me'n vaig");
         n_consumicions := 0;
      else
         -- Si n'hi ha, inicialitzem el nombre aleatori de consumicions
         n_consumicions := aleatori.Random(Gen);
         Put_Line(nom& " diu: Hola, avui fare" & Integer'Image(n_consumicions)& " consumicions");
      end if;

      -- Realitzem les consumicions
      for i in  1..n_consumicions loop
         Monitor.Consumir(i,n_consumicions,nom);
         delay Duration(aleatori.random(Gen)/3);
      end loop;

      -- Quan hem acabat les consumicions, informem al monitor
      final := (queden = 1); -- Si queden es 1, vol dir que no queden mes clients
      Put_Line(nom& " acaba i se'n va, queden" & Integer'Image(queden-1)& " clients>>>");
      Monitor.fi_client(queden);
   end tasca_client;

   -- Cos de la tasca reposador
   task body tasca_reposador is
      fi_reposador : boolean;
   begin
      Put_Line("El reposador"& Integer'Image(n_reposador)& " comenca a treballar");
      delay Duration(aleatori.random(Gen));

      -- Si no hi ha clients, omplim la maquina i mostrem el missatge de sortida
      if clients = 0 then
         Monitor.Reposar(n_reposador, fi_reposador);
      else
         -- Si hi ha clients, entrem al bucle fins que no en quedin
         while not final loop
            -- Omplim la maquina
            Monitor.Reposar(n_reposador,fi_reposador);
            delay Duration(aleatori.random(Gen));

            -- Comprovem si hem de dir adeu
            if not fi_reposador and final then
               Monitor.Reposar(n_reposador,fi_reposador);
            end if;
         end loop;
      end if;

      Put_Line("El reposador"& Integer'Image(n_reposador)& " acaba i se'n va >>>");
   end tasca_reposador;

   -- Definicio dels noms dels clients
   type matriu_noms is array (1 .. 10) of cadena_variable;
   noms : matriu_noms := (new String'("Aina"), new String'("Bernat"),
      new String'("Bel"), new String'("Miquel"), new String'("Joan"),
      new String'("Pau"), new String'("Cristina"), new String'("Andreu"),
      new String'("Maria"), new String'("Marta"));

   n_clients, n_reposadors: Integer;
   -- Tipus per poder crear l'array de tasques
   type acces_tasca_client is access all tasca_client;
   type acces_tasca_reposador is access all tasca_reposador;
   -- Arrays de fils de longitud variable que instanciarem en tenir tant el nombre de clients
   -- com el nombre de reposadors
   type fils_clients is array (Integer range <>) of acces_tasca_client;
   type fils_reposadors is array (Integer range <>) of acces_tasca_reposador;

begin
   -- Inicialitzacio del generador de nombres aleatoris
   aleatori.Reset(Gen);
   n_clients := aleatori.Random(Gen);  -- Nombre aleatori de clients
   queden := n_clients;
   n_reposadors := aleatori.Random(Gen);  -- Nombre aleatori de reposadors

   put_line("Simulacio amb" & Integer'Image(n_clients) & " clients i" & Integer'Image(n_reposadors) & " reposadors");

   -- Declarem la longitud dels arrays
   declare
      fils_client : fils_clients(1 .. n_clients);
      fils_reposador : fils_reposadors(1 .. n_reposadors);
   begin
      -- Inicialitzacio del monitor
      Monitor.iniciar_maquina(n_clients);

      -- Inicialitzem les tasques
      for i in 1 .. n_reposadors loop
         fils_reposador(i) := new tasca_reposador(i, n_clients);
      end loop;

      for i in 1 .. n_clients loop
         fils_client(i) := new tasca_client(noms(i), n_reposadors);
      end loop;
   end;

end maquinarefrescs;
