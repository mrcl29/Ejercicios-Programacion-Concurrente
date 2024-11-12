with Text_Io;         use Text_Io;
with Ada.Numerics.Discrete_Random;
with Ada.Numerics.Float_Random;
with maquinarefrescs; use maquinarefrescs;

procedure exercici2 is

    ---------------- VARIABLES ----------------

    subtype rang is Integer range 0 .. 10;

    package aleatori is new Ada.Numerics.Discrete_Random (rang);

    Numero_Aleatori : aleatori.Generator;

    ProductorConsumidor : Monitor_Productor_Consumidor;

    clients_restants : integer;

    fi : boolean := False;

    --------------------------------

    ---------------- REPOSADOR ----------------

    -- TYPE
    task type tasca_reposador
      (id             : Integer;
       nombre_clients : integer);
    -- BODY
    task body tasca_reposador is
        fi_reposador : boolean;
    begin
        Put_Line
          ("El reposador" & Integer'Image (id) & " comenca a treballar");

        delay Duration (aleatori.random (Numero_Aleatori));

        -- Logica per al reposador segons el nombre de clients
        if nombre_clients = 0 then
            ProductorConsumidor.Reposar (id, fi_reposador);
        else
            -- Si hi ha clients, entrem al bucle fins que no en quedin
            while not fi loop
                -- Omplim la maquina
                ProductorConsumidor.Reposar (id, fi_reposador);
                delay Duration (aleatori.random (Numero_Aleatori));

                -- Comprovem si hem de dir adeu
                if not fi_reposador and fi then
                    ProductorConsumidor.Reposar (id, fi_reposador);
                end if;
            end loop;
        end if;

        Put_Line
          ("El reposador"
           & Integer'Image (id)
           & " acaba i se'n va >>>>>>>>>>");
    end tasca_reposador;

    type acces_tasca_reposador is access all tasca_reposador;

    type fils_reposadors is array (Integer range <>) of acces_tasca_reposador;

    --------------------------------

    ---------------- CLIENT ----------------

    -- Definicio de tipus i variables per als clients
    type cadena_variable is access all String;
    type matriu_noms is array (1 .. 10) of cadena_variable;
    noms : matriu_noms :=
      (new String'("Aina"),
       new String'("Bernat"),
       new String'("Bel"),
       new String'("Albert"),
       new String'("Joan"),
       new String'("Pau"),
       new String'("Cristina"),
       new String'("Andreu"),
       new String'("Maria"),
       new String'("Marta"));

    -- TYPE
    task type tasca_client
      (nom_client        : cadena_variable;
       nombre_reposadors : integer);
    -- BODY
    task body tasca_client is
        nom                 : constant String := nom_client.all;
        refrescs_a_consumir : integer;
    begin
        -- Comportament del client segons la presencia de reposadors
        if (nombre_reposadors = 0) then
            Put_Line
              (nom & " diu: No hi ha reposadors a la maquina, me'n vaig");
            refrescs_a_consumir := 0;
        else
            -- Si n'hi ha, inicialitzem el nombre aleatori de consumicions
            refrescs_a_consumir := aleatori.Random (Numero_Aleatori);

            Put_Line
              (nom
               & " diu: Hola, avui fare"
               & Integer'Image (refrescs_a_consumir)
               & " consumicions");
        end if;

        -- Proces de consum de refrescs
        for i in 1 .. refrescs_a_consumir loop
            ProductorConsumidor.Agafar (i, refrescs_a_consumir, nom);

            delay Duration (aleatori.random (Numero_Aleatori) / 3);
        end loop;

        -- Finalitzacio del client
        fi := (clients_restants = 1);

        Put_Line
          (nom
           & " acaba i se'n va, queden"
           & Integer'Image (clients_restants - 1)
           & " clients >>>>>>>>>>");

        ProductorConsumidor.fi_client (clients_restants);
    end tasca_client;

    type acces_tasca_client is access all tasca_client;

    type fils_clients is array (Integer range <>) of acces_tasca_client;

    --------------------------------

    ---------------- VARIABLES EXECUCIO ----------------
    n_reposadors, n_clients : Integer;
    --------------------------------

    ---------------- EXECUCIO ----------------
begin
    -- Configuracio inicial
    aleatori.Reset (Numero_Aleatori);

    -- Generem un nombre aleatori de reposadors
    n_reposadors := aleatori.Random (Numero_Aleatori);

    -- Generem un nombre aleatori de clients
    n_clients := aleatori.Random (Numero_Aleatori);

    clients_restants := n_clients;

    put_line
      ("Simulacio amb"
       & Integer'Image (n_clients)
       & " clients i"
       & Integer'Image (n_reposadors)
       & " Reposadors");

    -- Execucio principal
    declare
        fils_reposador : fils_reposadors (1 .. n_reposadors);
        fils_client    : fils_clients (1 .. n_clients);
    begin
        -- Inicialitzacio del monitor
        ProductorConsumidor.iniciar (n_clients);

        -- Inicialitzem les tasques dels reposadors
        for i in 1 .. n_reposadors loop
            fils_reposador (i) := new tasca_reposador (i, n_clients);
        end loop;

        -- Inicialitzem les tasques dels clients
        for i in 1 .. n_clients loop
            fils_client (i) := new tasca_client (noms (i), n_reposadors);
        end loop;

        -- Esperem que totes les tasques dels reposadors acabin
        for i in 1 .. n_reposadors loop
            while not fils_reposador (i)'Terminated loop
                delay 0.1;
            end loop;
        end loop;

        -- Esperem que totes les tasques dels clients acabin
        for i in 1 .. n_clients loop
            while not fils_client (i)'Terminated loop
                delay 0.1;
            end loop;
        end loop;

        -- Esperem que l'usuari premi una tecla abans de tancar
        New_Line;
        Put_Line ("Presiona ENTER per sortir...");
        declare
            Aux : String := Get_Line;
        begin
            null;
        end;
    end;

    --------------------------------

end exercici2;
